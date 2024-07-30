package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
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
	navigationService := NewNavigationService()

	return navigationService.DetermineNextLocations(ctx, team)
}

func (s *GameplayService) CheckIn(ctx context.Context, team *models.Team, locationCode string) (response *ServiceResponse) {
	response = &ServiceResponse{}

	location, err := models.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Location code not found")
		response.AddFlashMessage(msg)
		response.Error = fmt.Errorf("finding location: %w", err)
		return response
	}

	if team.MustScanOut != "" {
		if locationCode != team.MustScanOut {
			msg := flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).SetTitle("You are already scanned in.")
			response.AddFlashMessage(msg)
			response.Error = fmt.Errorf("player must scan out first")
			return response
		}
	}

	// Check if the team has already scanned in
	scanned := false
	for _, s := range team.Scans {
		if s.LocationID == location.ID {
			scanned = true
			break
		}
	}
	if scanned {
		msg := flash.NewWarning("If you want to revisit this site, please click \"See my scanned locations\" below").SetTitle("You have already visited here.")
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
	_, err = location.LogScan(ctx, team.Code)
	if err != nil {
		msg := flash.NewWarning("Please double check the code and try again.").SetTitle("Couldn't scan in.")
		response.AddFlashMessage(msg)
		err := fmt.Errorf("logging scan: %w", err)
		response.Error = fmt.Errorf("logging scan: %w", err)
		return response
	}

	response.Data = make(map[string]interface{})
	response.Data["locationID"] = location.ID

	msg := flash.NewSuccess("You have scanned in.")
	response.AddFlashMessage(*msg)
	return response
}

func (s *GameplayService) CheckOut(ctx context.Context, teamCode, locationCode string) error {
	team, err := models.FindTeamByCode(ctx, teamCode)
	if err != nil {
		return fmt.Errorf("LogScanOut: %v", err)
	}

	location, err := models.FindLocationByInstanceAndCode(ctx, team.InstanceID, locationCode)
	if err != nil {
		return fmt.Errorf("LogScanOut: %v", err)
	}

	// Check if the team must scan out
	if team.MustScanOut == "" {
		return fmt.Errorf("You don't need to scan out.")
	} else if team.MustScanOut != locationCode {
		return fmt.Errorf("You need to scan out at %s", team.BlockingLocation.Name)
	}

	// Log the scan out
	err = location.LogScanOut(ctx, teamCode)
	if err != nil {
		return fmt.Errorf("LogScanOut: %v", err)
	}

	// Clear the mustScanOut field
	team.MustScanOut = ""
	err = team.Update(ctx)
	if err != nil {
		return fmt.Errorf("LogScanOut: %v", err)
	}

	return nil
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
