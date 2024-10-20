package services

import (
	"context"
	"errors"

	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type TeamService interface {
	// Update updates a team in the database
	Update(ctx context.Context, team *models.Team) error
	// Delete removes a team from the database
	Delete(ctx context.Context, teamCode string) error
	// AddTeams adds teams to the database
	AddTeams(ctx context.Context, instanceID string, count int) error
	// AwardPoints awards points to a team
	AwardPoints(ctx context.Context, team *models.Team, points int, reason string) error
	// LoadRelations loads relations for a team
	LoadRelation(ctx context.Context, team *models.Team, relation string) error
	// LoadRelations loads all relations for a team
	LoadRelations(ctx context.Context, team *models.Team) error
}

type teamService struct {
	teamRepo  repositories.TeamRepository
	batchSize int
}

// NewTeamService creates a new TeamService
func NewTeamService(tr repositories.TeamRepository) TeamService {
	return &teamService{
		teamRepo:  tr,
		batchSize: 100,
	}
}

// Update updates a team in the database
func (s *teamService) Update(ctx context.Context, team *models.Team) error {
	return s.teamRepo.Update(ctx, team)
}

// Delete removes a team from the database
func (s *teamService) Delete(ctx context.Context, teamCode string) error {
	return s.teamRepo.Delete(ctx, teamCode)
}

// AddTeams generates and inserts teams in batches, retrying if unique constraint errors occur
func (s *teamService) AddTeams(ctx context.Context, instanceID string, count int) error {
	for i := 0; i < count; i += s.batchSize {
		size := min(s.batchSize, count-i)
		teams := make([]models.Team, 0, size)

		for j := 0; j < size; j++ {
			var team models.Team
			for {
				// TODO: Remove magic number
				code := helpers.NewCode(4)
				team = models.Team{
					Code:       code,
					InstanceID: instanceID,
				}

				// Ensure code uniqueness within the current batch
				if !containsCode(teams, code) {
					teams = append(teams, team)
					break
				}
			}
		}

		// Insert the batch and retry if there's a unique constraint error
		err := s.teamRepo.InsertBatch(ctx, teams)
		if err != nil {
			if errors.Is(err, errors.New("unique constraint error")) {
				i -= s.batchSize // Retry this batch
				continue
			}
			return err // Return on other errors
		}
	}

	return nil
}

// Helper function to check for code uniqueness within a batch
func containsCode(teams []models.Team, code string) bool {
	for _, team := range teams {
		if team.Code == code {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AwardPoints awards points to a team
func (s *teamService) AwardPoints(ctx context.Context, team *models.Team, points int, _ string) error {
	team.Points += points
	return s.teamRepo.Update(ctx, team)
}

// LoadRelation loads the specified relation for a team
func (s *teamService) LoadRelation(ctx context.Context, team *models.Team, relation string) error {
	switch relation {
	case "Instance":
		return s.teamRepo.LoadInstance(ctx, team)
	case "Scans":
		return s.teamRepo.LoadScans(ctx, team)
	case "BlockingLocation":
		return s.teamRepo.LoadBlockingLocation(ctx, team)
	default:
		return errors.New("unknown relation")
	}
}

// LoadRelations loads all relations for a team
func (s *teamService) LoadRelations(ctx context.Context, team *models.Team) error {
	err := s.teamRepo.LoadRelations(ctx, team)
	if err != nil {
		return err
	}

	return nil
}
