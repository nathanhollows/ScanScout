package services

import (
	"context"
	"errors"
	"github.com/nathanhollows/Rapua/internal/models"
)

type PlayerGameService struct{}

// ScanLocation handles the logic for scanning a location
func (s *PlayerGameService) ScanLocation(ctx context.Context, team *models.Team, code string) error {
	location, err := models.FindLocationByCode(ctx, code)
	if err != nil {
		return errors.New("location not found")
	}

	if team.MustScanOut != "" {
		return errors.New("you must scan out before scanning a new location")
	}

	err = location.Marker.LogScan(ctx, team.Code)
	if err != nil {
		return err
	}

	// Additional game logic for scanning in (e.g., awarding points)
	return nil
}

// ScanOutLocation handles the logic for scanning out of a location
func (s *PlayerGameService) ScanOutLocation(ctx context.Context, team *models.Team, code string) error {
	location, err := models.FindLocationByCode(ctx, code)
	if err != nil {
		return errors.New("location not found")
	}

	if team.MustScanOut == "" || team.MustScanOut != code {
		return errors.New("you are not scanned in at this location")
	}

	err = location.Marker.LogScanOut(ctx, team.Code)
	if err != nil {
		return err
	}

	team.MustScanOut = ""
	return team.Update(ctx)
}

// GetNextLocations suggests the next locations for the team based on the navigation mode
// The navigation mode can be "free roam", "ordered", or "pseudo random"
func (s *PlayerGameService) GetNextLocations(ctx context.Context, team *models.Team) (*models.Locations, error) {
	switch team.Instance.NavigationMode {
	case models.FreeRoamMode:
		locs, err := models.FindAllLocations(ctx)
		return &locs, err
	case models.OrderedMode:
		return models.FindOrderedLocations(ctx, team)
	case models.PseudoRandomMode:
		return models.FindPseudoRandomLocations(ctx, team)
	default:
		return nil, errors.New("invalid navigation mode")
	}
}

// GetPlayerLocations returns the locations the player has scanned
func (s *PlayerGameService) GetPlayerLocations(ctx context.Context, team *models.Team) ([]*models.Location, error) {
	return team.GetVisitedLocations(ctx)
}

// GetLocationDetails returns the details of a specific location
func (s *PlayerGameService) GetLocationDetails(ctx context.Context, team *models.Team, code string) (*models.Location, error) {
	return models.FindLocationByCode(ctx, code)
}

// ParticipateInBonusEvent allows a player to participate in a bonus event
func (s *PlayerGameService) ParticipateInBonusEvent(ctx context.Context, team *models.Team, eventID string) error {
	event, err := models.FindEventByID(ctx, eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Type != "bonus" {
		return errors.New("invalid event type")
	}

	// Check if the team has already participated
	// Implement logic to check participation

	// Award points
	// Implement logic to award points

	return nil
}

// ParticipateInHuntEvent allows a player to participate in a hunt event
func (s *PlayerGameService) ParticipateInHuntEvent(ctx context.Context, team *models.Team, eventID string) error {
	event, err := models.FindEventByID(ctx, eventID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Type != "hunt" {
		return errors.New("invalid event type")
	}

	// Implement hunt event logic

	return nil
}
