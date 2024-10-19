package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type LocationService interface {
	FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error)
	LoadCluesForLocation(ctx context.Context, location *models.Location) error
	LoadCluesForLocations(ctx context.Context, locations *models.Locations) error
	LogCheckIn(ctx context.Context, team *models.Team, location *models.Location, mustCheckOut bool, validationRequired bool) (models.Scan, error)
}

type locationService struct {
	locationRepo repositories.LocationRepository
	clueRepo     repositories.ClueRepository
}

// NewLocationService creates a new instance of LocationService
func NewLocationService(clueRepo repositories.ClueRepository) LocationService {
	return locationService{
		clueRepo:     clueRepo,
		locationRepo: repositories.NewLocationRepository(),
	}
}

// FindLocationByInstanceAndCode finds a location by instance and code
func (s locationService) FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error) {
	location, err := s.locationRepo.FindLocationByInstanceAndCode(ctx, instanceID, code)
	if err != nil {
		fmt.Errorf("finding location by instance and code: %v", err)
		return nil, err
	}
	return location, nil
}

// LoadCluesForLocation loads the clues for a specific location if they are not already loaded
func (s locationService) LoadCluesForLocation(ctx context.Context, location *models.Location) error {
	if len(location.Clues) == 0 {
		clues, err := s.clueRepo.FindCluesByLocation(ctx, location.ID)
		if err != nil {
			slog.Error("error loading clues for location", "locationID", location.ID, "err", err)
			return err
		}
		location.Clues = clues
	}
	return nil
}

// LoadCluesForLocations loads the clues for all given locations if they are not already loaded
func (s locationService) LoadCluesForLocations(ctx context.Context, locations *models.Locations) error {
	for i := range *locations {
		err := s.LoadCluesForLocation(ctx, &(*locations)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// LogCheckIn logs a check in for a team at a location
func (s locationService) LogCheckIn(ctx context.Context, team *models.Team, location *models.Location, mustCheckOut bool, validationRequired bool) (models.Scan, error) {
	scan, err := s.locationRepo.LogCheckIn(ctx, team, location, mustCheckOut, validationRequired)
	if err != nil {
		return models.Scan{}, fmt.Errorf("logging check in: %w", err)
	}

	return scan, nil
}
