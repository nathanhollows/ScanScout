package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/blocks"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type BlockRepository interface {
	// GetContentBlocksByLocationID fetches all content blocks for a location
	GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error)
	// GetContentBlockByID fetches a content block by its ID
	GetByID(ctx context.Context, contentBlockID string) (blocks.Block, error)
	// SaveContentBlock saves a content block to the database
	Save(ctx context.Context, contentBlock *blocks.Block) error
}

type blockRepository struct{}

func NewBlockRepository() BlockRepository {
	return &blockRepository{}
}

// GetByLocationID fetches all content blocks for a location
func (r *blockRepository) GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error) {
	contentBlocks := models.Blocks{}
	err := db.DB.NewSelect().
		Model(&contentBlocks).
		Where("location_id = ?", locationID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return blocks.NewBlocksFromModel(contentBlocks)
}

// GetByID fetches a content block by its ID
func (r *blockRepository) GetByID(ctx context.Context, contentBlockID string) (blocks.Block, error) {
	contentBlock := &models.Block{}
	err := db.DB.NewSelect().
		Model(contentBlock).
		Where("id = ?", contentBlockID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return blocks.NewBlockFromModel(contentBlock)
}

// Save saves a content block to the database
func (r *blockRepository) Save(ctx context.Context, block *blocks.Block) error {
	contentBlock := models.Block{
		ID:   (*block).GetID(),
		Type: (*block).GetType(),
		Data: (*block).Data(),
	}
	if contentBlock.ID == "" {
		uuid := uuid.New()
		contentBlock.ID = uuid.String()
		_, err := db.DB.NewInsert().Model(&contentBlock).Exec(ctx)
		return err
	}
	_, err := db.DB.NewUpdate().Model(&contentBlock).WherePK().Exec(ctx)
	return err
}
