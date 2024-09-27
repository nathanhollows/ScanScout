package services

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/blocks"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type BlockService interface {
	// GetByLocationID fetches all content blocks for a location
	GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error)
	ValidateBlock(ctx context.Context, team models.Team, block blocks.Block) error
	ValidateBlocks(ctx context.Context, team models.Team, blocks blocks.Blocks) error
}

type blockService struct {
	Repo repositories.BlockRepository
}

func NewBlockService(repo repositories.BlockRepository) BlockService {
	return &blockService{
		Repo: repo,
	}
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
