package services

import (
	"context"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"log/slog"
)

type LocationService struct {
	clueRepo repositories.ClueRepository
}

// NewLocationService creates a new instance of LocationService
func NewLocationService(clueRepo repositories.ClueRepository) LocationService {
	return LocationService{
		clueRepo: clueRepo,
	}
}

// LoadCluesForLocation loads the clues for a specific location if they are not already loaded
func (s *LocationService) LoadCluesForLocation(ctx context.Context, location *models.Location) error {
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
func (s *LocationService) LoadCluesForLocations(ctx context.Context, locations *models.Locations) error {
	for i := range *locations {
		err := s.LoadCluesForLocation(ctx, &(*locations)[i])
		if err != nil {
			return err
		}
	}
	return nil
}
