package services

import (
	"context"
	"fmt"

	internalModels "github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/models"
)

type ClueService interface {
	UpdateClues(ctx context.Context, location *internalModels.Location, clues []string, clueIDs []string) error
}

type clueService struct {
	clueRepo     repositories.ClueRepository
	locationRepo repositories.LocationRepository
}

func NewClueService(clueRepo repositories.ClueRepository, locationRepo repositories.LocationRepository) ClueService {
	return &clueService{
		clueRepo:     clueRepo,
		locationRepo: locationRepo,
	}
}

func (s *clueService) UpdateClues(ctx context.Context, location *internalModels.Location, clues []string, clueIDs []string) error {
	var err error
	if len(location.Clues) == 0 {
		err = s.locationRepo.LoadClues(ctx, location)
	}
	if err != nil {
		return err
	}

	// There may be more clue IDs than clues, but not the other way around
	if len(clueIDs) > len(clues) {
		return fmt.Errorf("there are more clue IDs than clues")
	}

	// Delete all clues
	if len(location.Clues) > 0 {
		err = s.clueRepo.DeleteByLocationID(ctx, location.ID)
		if err != nil {
			return fmt.Errorf("deleting clues: %v", err)
		}
	}

	if len(clues) == 0 {
		return nil
	}

	// Add new clues, using the provided IDs if they exist
	for i, clue := range clues {
		if clue == "" {
			continue
		}
		c := &models.Clue{
			InstanceID: location.InstanceID,
			LocationID: location.ID,
			Content:    clue,
		}
		if i < len(clueIDs) {
			c.ID = clueIDs[i]
		}
		err = s.clueRepo.Save(ctx, c)
		if err != nil {
			return fmt.Errorf("saving clue: %v", err)
		}
	}

	return nil

}
