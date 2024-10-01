package services

import (
	"context"

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
	UpdateBlock(ctx context.Context, block *blocks.Block) error
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

// TODO: Finish
func (s *blockService) NewBlock(ctx context.Context, locationID string, blockType string) (blocks.Block, error) {
	block := s.createBlock(blockType)
	s.Repo.Create(ctx, &block, locationID)
	return block, nil
}

func (s *blockService) createBlock(blockType string) blocks.Block {
	switch blockType {
	case "markdown":
		return &blocks.MarkdownBlock{}
	case "password":
		return &blocks.PasswordBlock{}
	// case "checklist":
	// 	return &blocks.ChecklistBlock{}
	// case "api":
	// 	return &blocks.APIBlock{}
	default:
		return nil
	}
}

// UpdateBlock updates a block
func (s *blockService) UpdateBlock(ctx context.Context, block *blocks.Block) error {
	return s.Repo.Update(ctx, block)
}
