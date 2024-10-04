package repositories

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type TeamBlockStateRepository interface {
	GetByBlockAndTeam(ctx context.Context, blockID string, teamCode string) (models.TeamBlockState, error)
	Create(ctx context.Context, teamBlockState *models.TeamBlockState) error
	Save(ctx context.Context, teamBlockState *models.TeamBlockState) error
	Delete(ctx context.Context, block_id string, team_code string) error
}

type teamBlockStateRepository struct{}

func NewTeamBlockStateRepository() TeamBlockStateRepository {
	return &teamBlockStateRepository{}
}

// GetByBlockAndTeam fetches a specific team block state for a block and team
func (r *teamBlockStateRepository) GetByBlockAndTeam(ctx context.Context, blockID string, teamCode string) (models.TeamBlockState, error) {
	teamBlockState := models.TeamBlockState{}
	err := db.DB.NewSelect().
		Model(&teamBlockState).
		Where("block_id = ?", blockID).
		Where("team_code = ?", teamCode).
		Scan(ctx)
	if err != nil {
		return models.TeamBlockState{}, err
	}
	return teamBlockState, nil
}

// Create inserts a new team block state into the database
func (r *teamBlockStateRepository) Create(ctx context.Context, teamBlockState *models.TeamBlockState) error {
	_, err := db.DB.NewInsert().
		Model(teamBlockState).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Save modifies an existing team block state in the database
func (r *teamBlockStateRepository) Save(ctx context.Context, teamBlockState *models.TeamBlockState) error {
	_, err := db.DB.NewUpdate().
		Model(teamBlockState).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a team block state from the database by its ID
func (r *teamBlockStateRepository) Delete(ctx context.Context, block_id string, team_code string) error {
	_, err := db.DB.NewDelete().
		Model(&models.TeamBlockState{}).
		Where("block_id = ?", block_id).
		Where("team_code = ?", team_code).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
