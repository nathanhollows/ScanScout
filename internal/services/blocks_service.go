package services

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/blocks"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type BlockService interface {
	// GetByBlockID fetches a content block by its ID
	GetByBlockID(ctx context.Context, blockID string) (blocks.Block, error)
	// GetByLocationID fetches all content blocks for a location
	GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error)
	ValidateBlock(ctx context.Context, team models.Team, block blocks.Block) error
	ValidateBlocks(ctx context.Context, team models.Team, blocks blocks.Blocks) error
	NewBlock(ctx context.Context, locationID string, blockType string) (blocks.Block, error)
	UpdateBlock(ctx context.Context, block *blocks.Block, data map[string]string) error
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

// ValidateBlock validates a single block
func (s *blockService) ValidateBlock(ctx context.Context, team models.Team, block blocks.Block) error {
	return block.Validate(team.Code, nil)
}

// ValidateBlocks validates a slice of blocks
func (s *blockService) ValidateBlocks(ctx context.Context, team models.Team, blocks blocks.Blocks) error {
	return s.ValidateBlocks(ctx, team, blocks)
}

func (s *blockService) NewBlock(ctx context.Context, locationID string, blockType string) (blocks.Block, error) {
	var block blocks.Block
	switch blockType {
	case "markdown":
		block = &blocks.MarkdownBlock{}
	case "password":
		block = &blocks.PasswordBlock{}
	default:
		return nil, fmt.Errorf("block type %s not found", blockType)
	}
	s.Repo.Create(ctx, &block, locationID)
	return block, nil
}

// UpdateBlock updates a block
func (s *blockService) UpdateBlock(ctx context.Context, block *blocks.Block, data map[string]string) error {
	(*block).UpdateData(data)
	return s.Repo.Update(ctx, block)
}
