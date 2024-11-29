package services

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/models"
)

type BlockService interface {
	// GetByBlockID fetches a content block by its ID
	GetByBlockID(ctx context.Context, blockID string) (blocks.Block, error)
	// GetByLocationID fetches all content blocks for a location
	GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error)
	NewBlock(ctx context.Context, locationID string, blockType string) (blocks.Block, error)
	NewBlockState(ctx context.Context, blockID, teamCode string) (blocks.PlayerState, error)
	NewMockBlockState(ctx context.Context, blockID, teamCode string) (blocks.PlayerState, error)
	UpdateBlock(ctx context.Context, block blocks.Block, data map[string][]string) (blocks.Block, error)
	DeleteBlock(ctx context.Context, blockID string) error
	ReorderBlocks(ctx context.Context, locationID string, blockIDs []string) error
	GetBlocksWithStateByLocationIDAndTeamCode(ctx context.Context, locationID, teamCode string) ([]blocks.Block, map[string]blocks.PlayerState, error)
	GetBlockWithStateByBlockIDAndTeamCode(ctx context.Context, blockID, teamCode string) (blocks.Block, blocks.PlayerState, error)
	// CheckValidationRequiredForLocation checks if any blocks in a location require validation
	CheckValidationRequiredForLocation(ctx context.Context, locationID string) (bool, error)
	// CheckValidationRequiredForCheckIn checks if any blocks still require validation for a check in
	CheckValidationRequiredForCheckIn(ctx context.Context, locationID, teamCode string) (bool, error)
	// UpdateState updates the player state for a block
	UpdateState(ctx context.Context, state blocks.PlayerState) (blocks.PlayerState, error)
}

type blockService struct {
	blockRepo      repositories.BlockRepository
	blockStateRepo repositories.BlockStateRepository
}

func NewBlockService(blockRepo repositories.BlockRepository, blockStateRepo repositories.BlockStateRepository) BlockService {
	return &blockService{
		blockRepo:      blockRepo,
		blockStateRepo: blockStateRepo,
	}
}

// GetByBlockID fetches a content block by its ID
func (s *blockService) GetByBlockID(ctx context.Context, blockID string) (blocks.Block, error) {
	return s.blockRepo.GetByID(ctx, blockID)
}

// GetByLocationID fetches all content blocks for a location
func (s *blockService) GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error) {
	return s.blockRepo.GetByLocationID(ctx, locationID)
}

func (s *blockService) NewBlock(ctx context.Context, locationID string, blockType string) (blocks.Block, error) {
	// Use the blocks package to create the appropriate block based on the type.
	baseBlock := blocks.BaseBlock{
		Type:       blockType,
		LocationID: locationID,
	}

	// Let the blocks package handle the creation logic.
	block, err := blocks.CreateFromBaseBlock(baseBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to create block of type %s: %w", blockType, err)
	}

	// Store the new block in the repository.
	newBlock, err := s.blockRepo.Create(ctx, block, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to store block of type %s: %w", blockType, err)
	}

	return newBlock, nil
}

// NewBlockState creates a new block state
func (s *blockService) NewBlockState(ctx context.Context, blockID, teamCode string) (blocks.PlayerState, error) {
	state, err := s.blockStateRepo.NewBlockState(ctx, blockID, teamCode)
	if err != nil {
		return nil, fmt.Errorf("creating new block state: %w", err)
	}
	state, err = s.blockStateRepo.Create(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("storing new block state: %w", err)
	}
	return state, nil
}

// NewMockBlockState creates a new mock block state
func (s *blockService) NewMockBlockState(ctx context.Context, blockID, teamCode string) (blocks.PlayerState, error) {
	state, err := s.blockStateRepo.NewBlockState(ctx, blockID, teamCode)
	if err != nil {
		return nil, fmt.Errorf("creating new block state: %w", err)
	}
	return state, nil
}

// UpdateBlock updates a block
func (s *blockService) UpdateBlock(ctx context.Context, block blocks.Block, data map[string][]string) (blocks.Block, error) {
	err := block.UpdateBlockData(data)
	if err != nil {
		return nil, fmt.Errorf("updating block data: %w", err)
	}
	return s.blockRepo.Update(ctx, block)
}

// DeleteBlock deletes a block
func (s *blockService) DeleteBlock(ctx context.Context, blockID string) error {
	return s.blockRepo.Delete(ctx, blockID)
}

// ReorderBlocks reorders the blocks in a location
func (s *blockService) ReorderBlocks(ctx context.Context, locationID string, blockIDs []string) error {
	return s.blockRepo.Reorder(ctx, locationID, blockIDs)
}

func (s *blockService) GetBlocksWithStateByLocationIDAndTeamCode(ctx context.Context, locationID, teamCode string) ([]blocks.Block, map[string]blocks.PlayerState, error) {
	foundBlocks, states, err := s.blockRepo.GetBlocksAndStatesByLocationIDAndTeamCode(ctx, locationID, teamCode)
	if err != nil {
		return nil, nil, err
	}

	// Create a map for easier lookup of block states by block ID
	blockStates := make(map[string]blocks.PlayerState, len(states))
	for _, state := range states {
		blockStates[state.GetBlockID()] = state
	}

	return foundBlocks, blockStates, nil
}

func (s *blockService) GetBlockWithStateByBlockIDAndTeamCode(ctx context.Context, blockID, teamCode string) (blocks.Block, blocks.PlayerState, error) {
	if blockID == "" || teamCode == "" {
		return nil, nil, fmt.Errorf("blockID and teamCode must be set, got blockID: %s, teamCode: %s", blockID, teamCode)
	}

	return s.blockRepo.GetBlockAndStateByBlockIDAndTeamCode(ctx, blockID, teamCode)
}

// Convert block to model
func (s *blockService) ConvertBlockToModel(block blocks.Block) models.Block {
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

func (s *blockService) convertModelsToBlocks(cbs []models.Block) (blocks.Blocks, error) {
	b := make(blocks.Blocks, len(cbs))
	for i, cb := range cbs {
		block, err := s.convertModelToBlock(&cb)
		if err != nil {
			return nil, err
		}
		b[i] = block
	}
	return b, nil
}

func (s *blockService) convertModelToBlock(m *models.Block) (blocks.Block, error) {
	// Convert model to block
	newBlock, err := blocks.CreateFromBaseBlock(blocks.BaseBlock{
		ID:         m.ID,
		LocationID: m.LocationID,
		Type:       m.Type,
		Data:       m.Data,
		Order:      m.Ordering,
		Points:     m.Points,
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

// CheckValidationRequiredForLocation checks if any blocks in a location require validation
func (s *blockService) CheckValidationRequiredForLocation(ctx context.Context, locationID string) (bool, error) {
	blocks, err := s.GetByLocationID(ctx, locationID)
	if err != nil {
		return false, err
	}

	for _, block := range blocks {
		if block.RequiresValidation() {
			return true, nil
		}
	}

	return false, nil
}

// CheckValidationRequiredForCheckIn checks if any blocks still require validation for a check in
func (s *blockService) CheckValidationRequiredForCheckIn(ctx context.Context, locationID, teamCode string) (bool, error) {
	blocks, state, err := s.GetBlocksWithStateByLocationIDAndTeamCode(ctx, locationID, teamCode)
	if err != nil {
		return false, err
	}

	for _, block := range blocks {
		if block.RequiresValidation() {
			if state[block.GetID()] == nil {
				return true, nil
			}
			if state[block.GetID()].IsComplete() {
				continue
			}
			return true, nil
		}
	}

	return false, nil
}

// UpdateState updates the player state for a block
func (s *blockService) UpdateState(ctx context.Context, state blocks.PlayerState) (blocks.PlayerState, error) {
	return s.blockStateRepo.Update(ctx, state)
}
