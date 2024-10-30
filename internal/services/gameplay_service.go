package services

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"golang.org/x/exp/rand"
)

// Define errors
var (
	ErrTeamNotFound             = errors.New("team not found")
	ErrLocationNotFound         = errors.New("location not found")
	ErrTeamNotAllowedToCheckOut = errors.New("team not allowed to check out")
	ErrUnfinishedCheckIn        = errors.New("unfinished check in")
	ErrAlreadyCheckedIn         = errors.New("player has already scanned in")
	ErrUnecessaryCheckOut       = errors.New("player does not need to scan out")
	ErrInstanceSettingsNotFound = errors.New("instance settings not found")
)

type GameplayService interface {
	CheckGameStatus(ctx context.Context, team *models.Team) *ServiceResponse
	GetTeamByCode(ctx context.Context, teamCode string) (*models.Team, error)
	GetMarkerByCode(ctx context.Context, locationCode string) *ServiceResponse
	StartPlaying(ctx context.Context, teamCode, customTeamName string) *ServiceResponse
	SuggestNextLocations(ctx context.Context, team *models.Team, limit int) ([]models.Location, error)
	// CheckIn checks a team in at a location
	// It also manages the points and mustScanOut fields
	// As well as checking if any blocks must be completed
	CheckIn(ctx context.Context, team *models.Team, locationCode string) ServiceResponse
	CheckOut(ctx context.Context, team *models.Team, locationCode string) error
	CheckValidLocation(ctx context.Context, team *models.Team, locationCode string) (bool, error)
	ValidateAndUpdateBlockState(ctx context.Context, team models.Team, data map[string][]string) (blocks.PlayerState, blocks.Block, error)
}

type gameplayService struct {
	CheckInService    CheckInService
	LocationService   LocationService
	TeamService       TeamService
	BlockService      BlockService
	NavigationService NavigationService
	MarkerRepository  repositories.MarkerRepository
}

func NewGameplayService() GameplayService {
	return &gameplayService{
		CheckInService:  NewCheckInService(),
		LocationService: NewLocationService(repositories.NewClueRepository()),
		TeamService:     NewTeamService(repositories.NewTeamRepository()),
		BlockService: NewBlockService(
			repositories.NewBlockRepository(),
			repositories.NewBlockStateRepository(),
		),
		NavigationService: NewNavigationService(),
		MarkerRepository:  repositories.NewMarkerRepository(),
	}
}

// GetGameStatus returns the current status of the game
func (s *gameplayService) CheckGameStatus(ctx context.Context, team *models.Team) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	// Load the instance
	err := s.TeamService.LoadRelation(ctx, team, "Instance")
	if err != nil {
		response.Error = fmt.Errorf("loading instance: %w", err)
		return response
	}

	status := team.Instance.GetStatus()
	response.Data["status"] = status
	return response
}

func (s *gameplayService) GetTeamByCode(ctx context.Context, teamCode string) (*models.Team, error) {
	teamCode = strings.TrimSpace(strings.ToUpper(teamCode))
	team, err := s.TeamService.FindTeamByCode(ctx, teamCode)
	if err != nil {
		return nil, fmt.Errorf("GetTeamStatus: %w", err)
	}
	return team, nil
}

func (s *gameplayService) GetMarkerByCode(ctx context.Context, locationCode string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	locationCode = strings.TrimSpace(strings.ToUpper(locationCode))
	marker, err := s.MarkerRepository.FindByCode(ctx, locationCode)
	if err != nil {
		response.Error = fmt.Errorf("GetLocationByCode finding marker: %w", err)
		return response
	}
	response.Data["marker"] = marker
	return response
}

func (s *gameplayService) StartPlaying(ctx context.Context, teamCode, customTeamName string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	team, err := s.TeamService.FindTeamByCode(ctx, teamCode)
	if err != nil {
		response.Error = fmt.Errorf("StartPlaying find team: %w", err)
		response.AddFlashMessage(*flash.NewError("Team not found. Please double check the code and try again."))
		return response
	}

	// Update team with custom name if provided
	if !team.HasStarted || customTeamName != "" {
		team.Name = customTeamName
		team.HasStarted = true
		err = s.TeamService.Update(ctx, team)
		if err != nil {
			response.Error = fmt.Errorf("StartPlaying update team: %w", err)
			response.AddFlashMessage(*flash.NewError("Something went wrong. Please try again."))
			return response
		}
	}

	response.Data["team"] = team
	response.AddFlashMessage(*flash.NewSuccess("You have started the game!"))
	return response
}

func (s *gameplayService) SuggestNextLocations(ctx context.Context, team *models.Team, limit int) ([]models.Location, error) {
	// Populate the team with the necessary data
	err := s.TeamService.LoadRelations(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("loading relations: %w", err)
	}

	// Suggest the next locations for the team
	locations, err := s.NavigationService.DetermineNextLocations(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("determining next locations: %w", err)
	}

	return locations, nil
}

func (s *gameplayService) CheckIn(ctx context.Context, team *models.Team, locationCode string) (response ServiceResponse) {
	response = ServiceResponse{}

	// Load team relations
	err := s.TeamService.LoadRelations(ctx, team)
	if err != nil {
		response.Error = fmt.Errorf("loading relations: %w", err)
		return response
	}

	// A team may not check in if they must check out at a different location
	if team.MustCheckOut != "" && locationCode != team.MustCheckOut {
		response.Error = ErrAlreadyCheckedIn
		return response
	}

	// Find the location
	location, err := s.LocationService.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Location code not found")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("finding location: %w", err)
		return response
	}

	// A team may not check in if they have previously checked in at this location
	scanned := false
	for _, s := range team.CheckIns {
		if s.LocationID == location.ID {
			scanned = true
			break
		}
	}
	if scanned {
		response.Error = ErrAlreadyCheckedIn
		return response
	}

	// The location must be valid for the team to check in
	valid, err := s.CheckValidLocation(ctx, team, locationCode)
	if err != nil {
		response.Error = fmt.Errorf("check valid location: %w", err)
		return response
	}
	if !valid {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Location code not found")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("location not valid for team")
		return response
	}

	// Check if any blocks require validation (e.g. a checklist)
	validationRequired, err := s.BlockService.CheckValidationRequiredForLocation(ctx, location.ID)
	if err != nil {
		msg := flash.NewError("Something went wrong. Please try again.")
		response.AddFlashMessage(*msg)
		err := fmt.Errorf("checking if validation is required: %w", err)
		response.Error = fmt.Errorf("checking if validation is required: %w", err)
		return response
	}

	// Log the check in
	mustCheckOut := team.Instance.Settings.CompletionMethod == models.CheckInAndOut
	_, err = s.CheckInService.CheckIn(ctx, *team, *location, mustCheckOut, validationRequired)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Couldn't scan in.")
		response.AddFlashMessage(msg)
		err := fmt.Errorf("logging scan: %w", err)
		response.Error = fmt.Errorf("logging scan: %w", err)
		return response
	}

	err = s.LocationService.IncrementVisitorStats(ctx, location)
	if err != nil {
		msg := flash.NewError("Something went wrong. Please try again.")
		response.AddFlashMessage(*msg)
		response.Error = fmt.Errorf("incrementing visitor stats: %w", err)
		return response
	}

	// Points are only added if the team does not need to scan out
	// If the team must check out, the location is saved to the team
	if mustCheckOut {
		team.MustCheckOut = location.ID
	} else {
		team.Points += location.Points
	}
	err = s.TeamService.Update(ctx, team)
	if err != nil {
		msg := flash.NewError("Something went wrong. Please try again.")
		response.AddFlashMessage(*msg)
		response.Error = fmt.Errorf("updating team: %w", err)
		return response
	}

	response.Data = make(map[string]interface{})
	response.Data["location"] = location

	return response
}

func (s *gameplayService) CheckOut(ctx context.Context, team *models.Team, locationCode string) error {

	location, err := s.LocationService.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
	if err != nil {
		return fmt.Errorf("%w: finding location: %w", ErrLocationNotFound, err)
	}

	err = s.TeamService.LoadRelations(ctx, team)
	if err != nil {
		return fmt.Errorf("loading relations: %w", err)
	}

	// Check if the team must scan out
	if team.MustCheckOut == "" {
		return ErrUnecessaryCheckOut
	} else if team.MustCheckOut != location.ID {
		return ErrTeamNotAllowedToCheckOut
	}

	// Check if all blocks are completed
	unfinishedCheckIn, err := s.BlockService.CheckValidationRequiredForCheckIn(ctx, location.ID, team.Code)
	if err != nil {
		return fmt.Errorf("checking if validation is required: %w", err)
	}

	if unfinishedCheckIn {
		return ErrUnfinishedCheckIn
	}

	// Log the scan out
	_, err = s.CheckInService.CheckOut(ctx, team, location)
	if err != nil {
		return fmt.Errorf("logging scan out: %w", err)
	}

	return nil
}

// CheckLocation checks if the location is valid for the team to check in
func (s *gameplayService) CheckValidLocation(ctx context.Context, team *models.Team, locationCode string) (bool, error) {
	if team.Instance.ID == "" {
		return false, ErrInstanceNotFound
	}

	if team.Instance.Settings.InstanceID == "" {
		return false, ErrInstanceSettingsNotFound
	}

	if len(team.Instance.Locations) == 0 {
		return false, ErrLocationNotFound
	}

	valid, err := s.NavigationService.CheckValidLocation(ctx, team, &team.Instance.Settings, locationCode)
	if err != nil {
		return false, fmt.Errorf("checking if location is valid: %w", err)
	}

	return valid, nil
}

// loadClues loads the clues for the current location
// By default, it will only show one clue per location
func (s *gameplayService) loadClues(ctx context.Context, team *models.Team, locations []models.Location) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	err := s.LocationService.LoadCluesForLocations(ctx, &locations)
	if err != nil {
		response.Error = fmt.Errorf("loading clues for locations: %w", err)
		return response
	}

	// Create a seed for the random clue
	seed := team.Code
	h := fnv.New64a()
	_, err = h.Write([]byte(seed))
	if err != nil {
		response.Error = fmt.Errorf("creating seed for random clue: %w", err)
		return response
	}
	r := rand.New(rand.NewSource(uint64(h.Sum64())))

	// Randomly select a clue for each location
	// If the location has no clues, it will be skipped
	for i, location := range locations {
		if len(location.Clues) == 0 {
			continue
		} else if len(location.Clues) == 1 {
			continue
		}

		n := r.Intn(len(location.Clues))
		(locations)[i].Clues = location.Clues[n : n+1]
	}

	response.Data["nextLocations"] = locations
	return response
}

func (s *gameplayService) ValidateAndUpdateBlockState(ctx context.Context, team models.Team, data map[string][]string) (blocks.PlayerState, blocks.Block, error) {
	blockID := data["block"][0]
	if blockID == "" {
		return nil, nil, fmt.Errorf("blockID must be set")
	}

	block, state, err := s.BlockService.GetBlockWithStateByBlockIDAndTeamCode(ctx, blockID, team.Code)
	if err != nil {
		return nil, nil, fmt.Errorf("getting block with state: %w", err)
	}

	if state == nil {
		return nil, nil, fmt.Errorf("block state not found")
	}

	// Returning early here prevents the block from being updated
	// And points from being added to the team multiple times
	if state.IsComplete() {
		return state, block, nil
	}

	// Validate the block
	state, err = block.ValidatePlayerInput(state, data)
	if err != nil {
		return nil, nil, fmt.Errorf("validating block: %w", err)
	}

	state, err = s.BlockService.UpdateState(ctx, state)
	if err != nil {
		return nil, nil, fmt.Errorf("updating block state: %w", err)
	}

	// Assign points on completion
	if !state.IsComplete() {
		return state, block, nil
	}
	err = s.TeamService.AwardPoints(ctx, &team, block.GetPoints(), fmt.Sprint("Completed block ", block.GetName()))
	if err != nil {
		return nil, nil, fmt.Errorf("awarding points: %w", err)
	}

	// Update the check in all blocks have been completed
	unfinishedCheckIn, err := s.BlockService.CheckValidationRequiredForCheckIn(ctx, block.GetLocationID(), team.Code)
	if err != nil {
		return nil, nil, fmt.Errorf("checking if validation is required: %w", err)
	}

	if unfinishedCheckIn {
		return state, block, nil
	}

	err = s.CheckInService.CompleteBlocks(ctx, team.Code, block.GetLocationID())
	if err != nil {
		return nil, nil, fmt.Errorf("completing blocks: %w", err)
	}

	return state, block, nil
}
