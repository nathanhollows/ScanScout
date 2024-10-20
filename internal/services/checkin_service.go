package services

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type CheckInService interface {
	// CompleteBlocks marks all blocks for a location as complete
	CompleteBlocks(ctx context.Context, teamCode string, locationID string) error
	// LogCheckIn logs a check in for a team at a location
	LogCheckIn(ctx context.Context, team models.Team, location models.Location, mustCheckOut bool, validationRequired bool) (models.CheckIn, error)
}

type checkInService struct {
	checkInRepo repositories.CheckInRepository
}

func NewCheckInService() CheckInService {
	return &checkInService{
		checkInRepo: repositories.NewCheckInRepository(),
	}
}

func (s *checkInService) CompleteBlocks(ctx context.Context, teamCode string, locationID string) error {
	checkIn, err := s.checkInRepo.FindCheckInByTeamAndLocation(ctx, teamCode, locationID)
	if err != nil {
		return fmt.Errorf("finding check in: %w", err)
	}

	// If the check in is already complete, return early
	if checkIn.BlocksCompleted {
		return nil
	}

	checkIn.BlocksCompleted = true
	err = s.checkInRepo.Update(ctx, checkIn)
	if err != nil {
		return fmt.Errorf("updating check in: %w", err)
	}

	return nil

}

// LogCheckIn logs a check in for a team at a location
func (s *checkInService) LogCheckIn(ctx context.Context, team models.Team, location models.Location, mustCheckOut bool, validationRequired bool) (models.CheckIn, error) {
	scan, err := s.checkInRepo.LogCheckIn(ctx, team, location, mustCheckOut, validationRequired)
	if err != nil {
		return models.CheckIn{}, fmt.Errorf("logging check in: %w", err)
	}
	return scan, nil
}
