package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
)

type BlockStateRepository interface {
	GetByBlockAndTeam(ctx context.Context, blockID string, teamCode string) (blocks.PlayerState, error)
	Create(ctx context.Context, state blocks.PlayerState) (blocks.PlayerState, error)
	Update(ctx context.Context, block blocks.PlayerState) (blocks.PlayerState, error)
	Delete(ctx context.Context, block_id string, team_code string) error
	NewBlockState(ctx context.Context, blockID, teamCode string) (blocks.PlayerState, error)
}

type blockStateRepository struct{}

func NewBlockStateRepository() BlockStateRepository {
	return &blockStateRepository{}
}

// PlayerStateData struct implements the PlayerState interface
type PlayerStateData struct {
	blockID       string
	playerID      string
	playerData    json.RawMessage
	isComplete    bool
	pointsAwarded int
}

func (p *PlayerStateData) GetBlockID() string {
	return p.blockID
}

func (p *PlayerStateData) GetPlayerData() json.RawMessage {
	return p.playerData
}

func (p *PlayerStateData) GetPlayerID() string {
	return p.playerID
}

func (p *PlayerStateData) SetPlayerData(data json.RawMessage) {
	p.playerData = data
}

func (p *PlayerStateData) IsComplete() bool {
	return p.isComplete
}

func (p *PlayerStateData) SetComplete(complete bool) {
	p.isComplete = complete
}

func (p *PlayerStateData) GetPointsAwarded() int {
	return p.pointsAwarded
}

func (p *PlayerStateData) SetPointsAwarded(points int) {
	p.pointsAwarded = points
}

// Convert model state to PlayerState
func convertModelToPlayerStateData(state models.TeamBlockState) blocks.PlayerState {
	return &PlayerStateData{
		blockID:       state.BlockID,
		playerID:      state.TeamCode,
		playerData:    state.PlayerData,
		isComplete:    state.IsComplete,
		pointsAwarded: state.PointsAwarded,
	}
}

// Convert PlayerState to model state
func convertPlayerStateToModelData(state blocks.PlayerState) models.TeamBlockState {
	return models.TeamBlockState{
		TeamCode:      state.GetPlayerID(),
		BlockID:       state.GetBlockID(),
		PlayerData:    state.GetPlayerData(),
		IsComplete:    state.IsComplete(),
		PointsAwarded: state.GetPointsAwarded(),
	}
}

// GetByBlockAndTeam fetches a specific team block state for a block and team
func (r *blockStateRepository) GetByBlockAndTeam(ctx context.Context, blockID string, teamCode string) (blocks.PlayerState, error) {
	var modelState models.TeamBlockState
	err := db.DB.NewSelect().
		Model(&modelState).
		Where("block_id = ?", blockID).
		Where("team_code = ?", teamCode).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return convertModelToPlayerStateData(modelState), nil
}

// Create inserts a new team block state into the database
func (r *blockStateRepository) Create(ctx context.Context, state blocks.PlayerState) (blocks.PlayerState, error) {
	modelState := convertPlayerStateToModelData(state)

	_, err := db.DB.NewInsert().
		Model(&modelState).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return convertModelToPlayerStateData(modelState), nil
}

// Update modifies an existing team block state in the database
func (r *blockStateRepository) Update(ctx context.Context, state blocks.PlayerState) (blocks.PlayerState, error) {
	modelState := convertPlayerStateToModelData(state)
	if state.GetBlockID() == "" || state.GetPlayerID() == "" {
		return nil, fmt.Errorf("block_id and team_code must be set")
	}

	_, err := db.DB.NewUpdate().
		Model(&modelState).
		Set("player_data = ?", modelState.PlayerData).
		Set("is_complete = ?", modelState.IsComplete).
		Set("points_awarded = ?", modelState.PointsAwarded).
		Set("updated_at = ?", time.Now()).
		Where("block_id = ?", state.GetBlockID()).
		Where("team_code = ?", state.GetPlayerID()).
		Exec(ctx)

	return state, err
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

// NewBlockState creates a new block state
func (r *blockStateRepository) NewBlockState(ctx context.Context, blockID, teamCode string) (blocks.PlayerState, error) {
	state := &PlayerStateData{
		blockID:       blockID,
		playerID:      teamCode,
		playerData:    nil,
		isComplete:    false,
		pointsAwarded: 0,
	}
	return r.Create(ctx, state)
}
