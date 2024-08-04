package services

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"golang.org/x/exp/rand"
)

type GameplayService struct{}

func (s *GameplayService) GetTeamByCode(ctx context.Context, teamCode string) (*models.Team, error) {
	teamCode = strings.TrimSpace(strings.ToUpper(teamCode))
	team, err := models.FindTeamByCode(ctx, teamCode)
	if err != nil {
		return nil, fmt.Errorf("GetTeamStatus: %w", err)
	}
	return team, nil
}

func (s *GameplayService) GetMarkerByCode(ctx context.Context, locationCode string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	locationCode = strings.TrimSpace(strings.ToUpper(locationCode))
	marker, err := models.FindMarkerByCode(ctx, locationCode)
	if err != nil {
		response.Error = fmt.Errorf("GetLocationByCode finding marker: %w", err)
		return response
	}
	response.Data["marker"] = marker
	return response
}

func (s *GameplayService) StartPlaying(ctx context.Context, teamCode, customTeamName string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	response.Data = make(map[string]interface{})

	team, err := models.FindTeamByCode(ctx, teamCode)
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

func (s *GameplayService) SuggestNextLocations(ctx context.Context, team *models.Team, limit int) ServiceResponse {
	// Suggest the next locations for the team
	navigationService := NewNavigationService()
	response := navigationService.DetermineNextLocations(ctx, team)
	if response.Error != nil {
		response.Error = fmt.Errorf("suggesting next locations: %w", response.Error)
	}

	// Load clues for each location if necessary
	if team.Instance.Settings.NavigationMethod == models.ShowClues {
		response = s.loadClues(ctx, team, response.Data["nextLocations"].(models.Locations))
		if response.Error != nil {
			response.Error = fmt.Errorf("loading clues: %w", response.Error)
		}
	}

	return response
}

func (s *GameplayService) CheckIn(ctx context.Context, team *models.Team, locationCode string) (response ServiceResponse) {
	response = ServiceResponse{}

	location, err := models.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
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
			msg := flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).SetTitle("You are already scanned in.")
			response.AddFlashMessage(msg)
			response.Error = fmt.Errorf("player must scan out first")
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
		msg := flash.NewWarning("Try checking in somewhere else.").SetTitle("Repeat check-in")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("player has already scanned in")
		return response
	}

	// Check if the location is valid for the team to check in
	valid := s.CheckValidLocation(ctx, team, locationCode)
	if valid.Error != nil {
		response.Error = fmt.Errorf("check valid location: %w", valid.Error)
		return response
	}

	// Log the CheckIn
	mustCheckOut := team.Instance.Settings.CompletionMethod == models.CheckInAndOut
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

	msg := flash.NewSuccess("You have scanned in.")
	response.AddFlashMessage(*msg)
	return response
}

func (s *GameplayService) CheckOut(ctx context.Context, team *models.Team, locationCode string) (response ServiceResponse) {
	response = ServiceResponse{}

	location, err := models.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
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
func (s *GameplayService) CheckValidLocation(ctx context.Context, team *models.Team, locationCode string) (response *ServiceResponse) {
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
func (s *GameplayService) loadClues(ctx context.Context, team *models.Team, locations models.Locations) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	for i := range locations {
		(locations)[i].LoadClues(ctx)
	}

	// Create a seed for the random clue
	seed := team.Code
	h := fnv.New64a()
	_, err := h.Write([]byte(seed))
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
