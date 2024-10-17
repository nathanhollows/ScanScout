package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
)

type TeamRepository interface {
	// Update saves or updates a team in the database
	Update(ctx context.Context, t *models.Team) error
	// Delete removes the team from the database
	Delete(ctx context.Context, teamCode string) error
	// FindAll returns all teams for an instance
	FindAll(ctx context.Context, instanceID string) ([]models.Team, error)
	// FindAllWithScans returns all teams for an instance with scans
	FindAllWithScans(ctx context.Context, instanceID string) ([]models.Team, error)
	// AddTeams adds teams to the database
	InsertBatch(ctx context.Context, teams []models.Team) error
}

type teamRepository struct{}

// NewTeamRepository creates a new TeamRepository
func NewTeamRepository() TeamRepository {
	return &teamRepository{}
}

// Update saves or updates a team in the database
func (r *teamRepository) Update(ctx context.Context, t *models.Team) error {
	_, err := db.DB.NewUpdate().Model(t).WherePK().Exec(ctx)
	return err
}

func (r *teamRepository) Delete(ctx context.Context, teamCode string) error {
	_, err := db.DB.NewDelete().Model(&models.Team{Code: teamCode}).WherePK().Exec(ctx)
	return err
}

func (r *teamRepository) FindAll(ctx context.Context, instanceID string) ([]models.Team, error) {
	var teams []models.Team
	err := db.DB.NewSelect().
		Model(&teams).
		Where("team.instance_id = ?", instanceID).
		Scan(ctx)
	if err != nil {
		return teams, err
	}
	return teams, nil
}

func (r *teamRepository) FindAllWithScans(ctx context.Context, instanceID string) ([]models.Team, error) {
	var teams []models.Team
	err := db.DB.NewSelect().
		Model(&teams).
		Where("team.instance_id = ?", instanceID).
		// Add the scans in the relation order by location_id
		Relation("Scans", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("location_id ASC")
		}).
		Scan(ctx)
	if err != nil {
		return teams, err
	}
	return teams, nil
}

// InsertBatch inserts a batch of teams and returns an error if there's a unique constraint conflict
func (r *teamRepository) InsertBatch(ctx context.Context, teams []models.Team) error {
	_, err := db.DB.NewInsert().Model(&teams).Exec(ctx)
	if err != nil && isUniqueConstraintError(err) {
		return errors.New("unique constraint error")
	}
	return err
}

// isUniqueConstraintError checks if an error is due to a unique constraint violation
func isUniqueConstraintError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unique constraint")
}
