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
	UpdateBlock(ctx context.Context, block *blocks.Block, data map[string]string) error
	DeleteBlock(ctx context.Context, blockID string) error
	ReorderBlocks(ctx context.Context, locationID string, blockIDs []string) error
	GetBlocksWithStateByLocationIDAndTeamCode(ctx context.Context, locationID, teamCode string) ([]blocks.Block, map[string]models.TeamBlockState, error)
	GetBlockWithStateByBlockIDAndTeamCode(ctx context.Context, blockID, teamCode string) (blocks.Block, models.TeamBlockState, error)
}

type blockService struct {
	Repo repositories.BlockRepository
}

func NewBlockService(repo repositories.BlockRepository) BlockService {
	return &blockService{
		Repo: repo,
	}
}

// GetByBlockID fetches a content block by its ID
func (s *blockService) GetByBlockID(ctx context.Context, blockID string) (blocks.Block, error) {
	return s.Repo.GetByID(ctx, blockID)
}

// GetByLocationID fetches all content blocks for a location
func (s *blockService) GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error) {
	return s.Repo.GetByLocationID(ctx, locationID)
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
	err = s.Repo.Create(ctx, &block, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to store block of type %s: %w", blockType, err)
	}

	return block, nil
}

// UpdateBlock updates a block
func (s *blockService) UpdateBlock(ctx context.Context, block *blocks.Block, data map[string]string) error {
	(*block).UpdateBlockData(data)
	return s.Repo.Update(ctx, block)
}

// DeleteBlock deletes a block
func (s *blockService) DeleteBlock(ctx context.Context, blockID string) error {
	return s.Repo.Delete(ctx, blockID)
}

// ReorderBlocks reorders the blocks in a location
func (s *blockService) ReorderBlocks(ctx context.Context, locationID string, blockIDs []string) error {
	return s.Repo.Reorder(ctx, locationID, blockIDs)
}

func (s *blockService) GetBlocksWithStateByLocationIDAndTeamCode(ctx context.Context, locationID, teamCode string) ([]blocks.Block, map[string]models.TeamBlockState, error) {
	blocks, states, err := s.Repo.GetBlocksAndStatesByLocationIDAndTeamCode(ctx, locationID, teamCode)
	if err != nil {
		return nil, nil, err
	}

	blockObjects, err := s.convertModelsToBlocks(blocks)
	if err != nil {
		return nil, nil, err
	}

	// Create a map for easier lookup of block states by block ID
	blockStates := make(map[string]models.TeamBlockState)
	for _, state := range states {
		blockStates[state.BlockID] = state
	}

	return blockObjects, blockStates, nil
}

func (s *blockService) GetBlockWithStateByBlockIDAndTeamCode(ctx context.Context, blockID, teamCode string) (blocks.Block, models.TeamBlockState, error) {
	blockModel, state, err := s.Repo.GetBlockAndStateByBlockIDAndTeamCode(ctx, blockID, teamCode)
	if err != nil {
		return nil, models.TeamBlockState{}, err
	}

	blockObject, err := s.convertModelToBlock(&blockModel)
	if err != nil {
		return nil, models.TeamBlockState{}, err
	}

	return blockObject, state, nil
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
