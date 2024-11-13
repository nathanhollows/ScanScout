package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
)

type BlockRepository interface {
	// GetBlocksByLocationID fetches all blocks for a location
	GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error)
	// GetBlockByID fetches a block by its ID
	GetByID(ctx context.Context, blockID string) (blocks.Block, error)
	// Saveblock saves a block to the database
	Save(ctx context.Context, block blocks.Block) (blocks.Block, error)
	Create(ctx context.Context, block blocks.Block, locationID string) (blocks.Block, error)
	Update(ctx context.Context, block blocks.Block) (blocks.Block, error)
	Delete(ctx context.Context, blockID string) error
	// Delete by location ID
	DeleteByLocationID(ctx context.Context, locationID string) error
	Reorder(ctx context.Context, locationID string, blockIDs []string) error
	GetBlocksAndStatesByLocationIDAndTeamCode(ctx context.Context, locationID string, teamCode string) ([]blocks.Block, []blocks.PlayerState, error)
	GetBlockAndStateByBlockIDAndTeamCode(ctx context.Context, blockID string, teamCode string) (blocks.Block, blocks.PlayerState, error)
}

type blockRepository struct{}

func NewBlockRepository() BlockRepository {
	return &blockRepository{}
}

// GetByLocationID fetches all blocks for a location
func (r *blockRepository) GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error) {
	modelBlocks := []models.Block{}
	err := db.DB.NewSelect().
		Model(&modelBlocks).
		Where("location_id = ?", locationID).
		Order("ordering ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return convertModelsToBlocks(modelBlocks)
}

// GetByID fetches a block by its ID
func (r *blockRepository) GetByID(ctx context.Context, blockID string) (blocks.Block, error) {
	modelBlock := &models.Block{}
	err := db.DB.NewSelect().
		Model(modelBlock).
		Where("id = ?", blockID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return convertModelToBlock(modelBlock)
}

// Save saves a block to the database
func (r *blockRepository) Save(ctx context.Context, block blocks.Block) (blocks.Block, error) {
	model := convertBlockToModel(block)
	if model.ID == "" {
		model.ID = uuid.New().String()
		_, err := db.DB.NewInsert().Model(&model).Exec(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := db.DB.NewUpdate().Model(&model).WherePK().Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// Convert back to block and return
	updatedBlock, err := convertModelToBlock(&model)
	if err != nil {
		return nil, err
	}
	return updatedBlock, nil
}

// Create saves a new block to the database
func (r *blockRepository) Create(ctx context.Context, block blocks.Block, locationID string) (blocks.Block, error) {
	modelBlock := models.Block{
		ID:                 uuid.New().String(),
		LocationID:         block.GetLocationID(),
		Type:               block.GetType(),
		Data:               block.GetData(),
		Ordering:           1e4,
		Points:             block.GetPoints(),
		ValidationRequired: block.RequiresValidation(),
	}
	_, err := db.DB.NewInsert().Model(&modelBlock).Exec(ctx)
	if err != nil {
		return nil, err
	}
	// Convert back to block and return
	createdBlock, err := convertModelToBlock(&modelBlock)
	if err != nil {
		return nil, err
	}
	return createdBlock, nil
}

// Update saves an existing block to the database
func (r *blockRepository) Update(ctx context.Context, block blocks.Block) (blocks.Block, error) {
	modelBlock := convertBlockToModel(block)
	_, err := db.DB.NewUpdate().Model(&modelBlock).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	// Convert back to block and return
	updatedBlock, err := convertModelToBlock(&modelBlock)
	if err != nil {
		return nil, err
	}
	return updatedBlock, nil
}

// Convert block to model
func convertBlockToModel(block blocks.Block) models.Block {
	return models.Block{
		ID:                 block.GetID(),
		LocationID:         block.GetLocationID(),
		Type:               block.GetType(),
		Ordering:           block.GetOrder(),
		Data:               block.GetData(),
		Points:             block.GetPoints(),
		ValidationRequired: block.RequiresValidation(),
	}
}

func convertModelsToBlocks(modelBlocks []models.Block) (blocks.Blocks, error) {
	b := make(blocks.Blocks, len(modelBlocks))
	for i, modelBlock := range modelBlocks {
		block, err := convertModelToBlock(&modelBlock)
		if err != nil {
			return nil, err
		}
		b[i] = block
	}
	return b, nil
}

func convertModelToBlock(model *models.Block) (blocks.Block, error) {
	// Convert model to block
	newBlock, err := blocks.CreateFromBaseBlock(blocks.BaseBlock{
		ID:         model.ID,
		LocationID: model.LocationID,
		Type:       model.Type,
		Data:       model.Data,
		Order:      model.Ordering,
		Points:     model.Points,
	})
	if err != nil {
		return nil, err
	}
	err = newBlock.ParseData()
	if err != nil {
		return nil, err
	}
	return newBlock, nil
}

// Delete deletes a block from the database
func (r *blockRepository) Delete(ctx context.Context, blockID string) error {
	_, err := db.DB.NewDelete().Model(&models.Block{}).Where("id = ?", blockID).Exec(ctx)
	return err
}

// DeleteByLocationID deletes all blocks for a location
func (r *blockRepository) DeleteByLocationID(ctx context.Context, locationID string) error {
	_, err := db.DB.NewDelete().Model(&models.Block{}).Where("location_id = ?", locationID).Exec(ctx)
	return err
}

// Reorder reorders the blocks
func (r *blockRepository) Reorder(ctx context.Context, locationID string, blockIDs []string) error {
	for i, blockID := range blockIDs {
		_, err := db.DB.NewUpdate().
			Model(&models.Block{}).
			Set("ordering = ?", i).
			Where("id = ?", blockID).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetBlocksAndStatesByLocationIDAndTeamCode fetches all blocks for a location with their player states
func (r *blockRepository) GetBlocksAndStatesByLocationIDAndTeamCode(ctx context.Context, locationID string, teamCode string) ([]blocks.Block, []blocks.PlayerState, error) {
	modelBlocks := []models.Block{}
	states := []models.TeamBlockState{}

	err := db.DB.NewSelect().
		Model(&modelBlocks).
		Where("location_id = ?", locationID).
		Order("ordering ASC").
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	if teamCode != "" {
		err = db.DB.NewSelect().
			Model(&states).
			Where("block_id IN (?)", db.DB.NewSelect().Model((*models.Block)(nil)).Column("id").Where("location_id = ?", locationID)).
			Where("team_code = ?", teamCode).
			Scan(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	foundBlocks, err := convertModelsToBlocks(modelBlocks)
	if err != nil {
		return nil, nil, err
	}

	playerStates := make([]blocks.PlayerState, len(states))
	for i, state := range states {
		playerStates[i] = convertModelToPlayerStateData(state)
	}

	// Populate playerStates with empty states for blocks without a state
	for _, block := range foundBlocks {
		found := false
		for _, state := range playerStates {
			if state.GetBlockID() == block.GetID() {
				found = true
				break
			}
		}
		if !found {
			playerStates = append(playerStates, &PlayerStateData{
				blockID:    block.GetID(),
				playerID:   "",
				playerData: []byte("{}"),
			})

		}
	}

	return foundBlocks, playerStates, nil
}

// GetBlockAndStateByBlockIDAndTeamCode fetches a block by its ID with the player state for a given team
func (r *blockRepository) GetBlockAndStateByBlockIDAndTeamCode(ctx context.Context, blockID string, teamCode string) (blocks.Block, blocks.PlayerState, error) {
	stateRepo := NewBlockStateRepository()

	modelBlock := models.Block{}
	err := db.DB.NewSelect().
		Model(&modelBlock).
		Where("id = ?", blockID).
		Scan(ctx)
	if err != nil {
		return nil, nil, err
	}

	state, err := stateRepo.GetByBlockAndTeam(ctx, blockID, teamCode)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, nil, err
	} else if err != nil {
		state, err = stateRepo.NewBlockState(ctx, blockID, teamCode)
		if err != nil {
			return nil, nil, err
		}
	}

	block, err := convertModelToBlock(&modelBlock)
	if err != nil {
		return nil, nil, err
	}

	return block, state, nil
}
