package services

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/internal/flash"
	internalModels "github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/models"
	"golang.org/x/exp/rand"
)

// Define errors
var (
	ErrTeamNotFound     = errors.New("team not found")
	ErrAlreadyCheckedIn = errors.New("player has already scanned in")
)

type GameplayService interface {
	CheckGameStatus(ctx context.Context, team *internalModels.Team) *ServiceResponse
	GetTeamByCode(ctx context.Context, teamCode string) (*internalModels.Team, error)
	GetMarkerByCode(ctx context.Context, locationCode string) *ServiceResponse
	StartPlaying(ctx context.Context, teamCode, customTeamName string) *ServiceResponse
	SuggestNextLocations(ctx context.Context, team *internalModels.Team, limit int) ServiceResponse
	CheckIn(ctx context.Context, team *internalModels.Team, locationCode string) ServiceResponse
	CheckOut(ctx context.Context, team *internalModels.Team, locationCode string) ServiceResponse
	CheckValidLocation(ctx context.Context, team *internalModels.Team, locationCode string) *ServiceResponse
	ValidateAndUpdateBlockState(ctx context.Context, block blocks.Block, state *models.TeamBlockState, data map[string]string) error
}

type gameplayService struct {
	LocationService LocationService
}

func NewGameplayService() GameplayService {
	return &gameplayService{
		LocationService: NewLocationService(repositories.NewClueRepository()),
	}
}

// GetGameStatus returns the current status of the game
func (s *gameplayService) CheckGameStatus(ctx context.Context, team *internalModels.Team) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	// Load the instance
	err := team.LoadInstance(ctx)
	if err != nil {
		response.Error = fmt.Errorf("loading instance: %w", err)
		return response
	}

	status := team.Instance.GetStatus()
	response.Data["status"] = status
	return response
}

func (s *gameplayService) GetTeamByCode(ctx context.Context, teamCode string) (*internalModels.Team, error) {
	teamCode = strings.TrimSpace(strings.ToUpper(teamCode))
	team, err := internalModels.FindTeamByCode(ctx, teamCode)
	if err != nil {
		return nil, fmt.Errorf("GetTeamStatus: %w", err)
	}
	return team, nil
}

func (s *gameplayService) GetMarkerByCode(ctx context.Context, locationCode string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	locationCode = strings.TrimSpace(strings.ToUpper(locationCode))
	marker, err := internalModels.FindMarkerByCode(ctx, locationCode)
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

	team, err := internalModels.FindTeamByCode(ctx, teamCode)
	if err != nil {
		response.Error = fmt.Errorf("StartPlaying find team: %w", err)
		response.AddFlashMessage(*flash.NewError("Team not found. Please double check the code and try again."))
		return response
	}

	// Update team with custom name if provided
	if !team.HasStarted || customTeamName != "" {
		team.Name = customTeamName
		team.HasStarted = true
		if err := team.Update(ctx); err != nil {
			response.Error = fmt.Errorf("StartPlaying update team: %w", err)
			response.AddFlashMessage(*flash.NewError("Something went wrong. Please try again."))
			return response
		}
	}

	response.Data["team"] = team
	response.AddFlashMessage(*flash.NewSuccess("You have started the game!"))
	return response
}

func (s *gameplayService) SuggestNextLocations(ctx context.Context, team *internalModels.Team, limit int) ServiceResponse {
	// Suggest the next locations for the team
	navigationService := NewNavigationService()
	response := navigationService.DetermineNextLocations(ctx, team)
	if response.Error != nil {
		response.Error = fmt.Errorf("suggesting next locations: %w", response.Error)
	}

	// Load clues for each location if necessary
	if team.Instance.Settings.NavigationMethod == internalModels.ShowClues {
		response = s.loadClues(ctx, team, response.Data["nextLocations"].(internalModels.Locations))
		if response.Error != nil {
			response.Error = fmt.Errorf("loading clues: %w", response.Error)
		}
	}

	return response
}

func (s *gameplayService) CheckIn(ctx context.Context, team *internalModels.Team, locationCode string) (response ServiceResponse) {
	response = ServiceResponse{}

	location, err := internalModels.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Location code not found")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("finding location: %w", err)
		return response
	}

	if team.MustScanOut != "" {
		if locationCode != team.MustScanOut {
			err := team.LoadBlockingLocation(ctx)
			if err != nil {
				response.Error = fmt.Errorf("loading blocking location: %w", err)
				return response
			}
			response.Error = ErrAlreadyCheckedIn
			return response
		}
	}

	// Check if the team has already scanned in
	scanned := false
	team.LoadScans(ctx)
	for _, s := range team.Scans {
		if s.LocationID == location.ID {
			scanned = true
			break
		}
	}
	if scanned {
		response.Error = ErrAlreadyCheckedIn
		return response
	}

	// Check if the location is valid for the team to check in
	valid := s.CheckValidLocation(ctx, team, locationCode)
	if valid.Error != nil {
		response.Error = fmt.Errorf("check valid location: %w", valid.Error)
		return response
	}

	// Log the CheckIn
	mustCheckOut := team.Instance.Settings.CompletionMethod == internalModels.CheckInAndOut
	_, err = location.LogCheckIn(ctx, *team, mustCheckOut)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Couldn't scan in.")
		response.AddFlashMessage(msg)
		err := fmt.Errorf("logging scan: %w", err)
		response.Error = fmt.Errorf("logging scan: %w", err)
		return response
	}

	response.Data = make(map[string]interface{})
	response.Data["location"] = location

	return response
}

func (s *gameplayService) CheckOut(ctx context.Context, team *internalModels.Team, locationCode string) (response ServiceResponse) {
	response = ServiceResponse{}

	location, err := internalModels.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Location code not found")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("finding location: %w", err)
		return response
	}

	err = team.LoadBlockingLocation(ctx)
	if err != nil {
		msg := flash.NewError("Something went wrong. Please try again.")
		response.AddFlashMessage(*msg)
		response.Error = fmt.Errorf("loading blocking location: %w", err)
		return response
	}

	// Check if the team must scan out
	if team.MustScanOut == "" {
		msg := flash.NewDefault("You don't need to scan out.")
		response.AddFlashMessage(*msg)
		return response
	} else if team.MustScanOut != location.ID {
		msg := flash.NewWarning("You need to scan out at " + team.BlockingLocation.Name + ".")
		response.AddFlashMessage(*msg)
		return response
	}

	// Log the scan out
	err = location.LogScanOut(ctx, team)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Couldn't scan out.")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("logging scan out: %w", err)
		return response
	}

	// Clear the mustScanOut field
	team.MustScanOut = ""
	err = team.Update(ctx)
	if err != nil {
		msg := flash.NewError("Something went wrong. Please try again.")
		response.AddFlashMessage(*msg)
		response.Error = fmt.Errorf("updating team: %w", err)
		return response
	}

	return response
}

// CheckLocation checks if the location is valid for the team to check in
func (s *gameplayService) CheckValidLocation(ctx context.Context, team *internalModels.Team, locationCode string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	NavigationService := NewNavigationService()

	err := team.LoadInstance(ctx)
	if err != nil {
		response.Error = fmt.Errorf("loading instance on team: %w", err)
		return response
	}
	err = team.Instance.LoadSettings(ctx)
	if err != nil {
		response.Error = fmt.Errorf("loading settings on instance: %w", err)
		return response
	}

	resp := NavigationService.CheckValidLocation(ctx, team, &team.Instance.Settings, locationCode)
	if resp.Error != nil {
		response.Error = fmt.Errorf("checking if code matches valid locations: %w", resp.Error)
		return response
	}

	return response
}

// loadClues loads the clues for the current location
// By default, it will only show one clue per location
func (s *gameplayService) loadClues(ctx context.Context, team *internalModels.Team, locations internalModels.Locations) (response ServiceResponse) {
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

func (s *gameplayService) ValidateAndUpdateBlockState(ctx context.Context, block blocks.Block, state *models.TeamBlockState, data map[string]string) error {
	blockStateRepo := repositories.NewBlockStateRepository()

	// Check if the block is already complete
	if state.IsComplete {
		return nil
	}

	// Validate the block
	gameErr := block.ValidatePlayerInput(state, data)

	err := blockStateRepo.Save(ctx, state)
	if err != nil {
		return fmt.Errorf("saving block state: %w", err)
	}

	if gameErr != nil {
		return fmt.Errorf("validating block: %w", gameErr)
	}

	return nil
}
