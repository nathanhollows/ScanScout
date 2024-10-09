package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type BlockStateRepository interface {
	GetByBlockAndTeam(ctx context.Context, blockID string, teamCode string) (models.TeamBlockState, error)
	Create(ctx context.Context, teamBlockState *models.TeamBlockState) error
	Save(ctx context.Context, teamBlockState *models.TeamBlockState) error
	Update(ctx context.Context, teamBlockState *models.TeamBlockState) error
	Delete(ctx context.Context, block_id string, team_code string) error
}

type blockStateRepository struct{}

func NewBlockStateRepository() BlockStateRepository {
	return &blockStateRepository{}
}

// GetByBlockAndTeam fetches a specific team block state for a block and team
func (r *blockStateRepository) GetByBlockAndTeam(ctx context.Context, blockID string, teamCode string) (models.TeamBlockState, error) {
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
func (r *blockStateRepository) Create(ctx context.Context, teamBlockState *models.TeamBlockState) error {
	_, err := db.DB.NewInsert().
		Model(teamBlockState).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Save modifies an existing team block state in the database
func (r *blockStateRepository) Save(ctx context.Context, teamBlockState *models.TeamBlockState) error {
	if teamBlockState.BlockID == "" || teamBlockState.TeamCode == "" {
		return fmt.Errorf("block_id and team_code must be set")
	}
	if teamBlockState.ID != "" {
		return r.Update(ctx, teamBlockState)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	teamBlockState.ID = id.String()

	_, err = db.DB.NewInsert().
		Model(teamBlockState).
		Exec(ctx)

	return err
}

// Update modifies an existing team block state in the database
func (r *blockStateRepository) Update(ctx context.Context, teamBlockState *models.TeamBlockState) error {
	if teamBlockState.ID == "" {
		return fmt.Errorf("id must be set")
	}
	if teamBlockState.BlockID == "" || teamBlockState.TeamCode == "" {
		return fmt.Errorf("block_id and team_code must be set")
	}
	_, err := db.DB.NewUpdate().
		Model(teamBlockState).
		WherePK("id").
		Exec(ctx)
	return err
}

// Delete removes a team block state from the database by its ID
func (r *blockStateRepository) Delete(ctx context.Context, block_id string, team_code string) error {
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
