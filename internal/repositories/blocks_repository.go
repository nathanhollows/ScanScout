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
	Create(ctx context.Context, contentBlock *blocks.Block, locationID string) error
	Update(ctx context.Context, contentBlock *blocks.Block) error
	Delete(ctx context.Context, contentBlockID string) error
	Reorder(ctx context.Context, locationID string, blockIDs []string) error
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
		Order("ordering ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertModelsToBlocks(contentBlocks)
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
	return ConvertModelToBlock(contentBlock)
}

// Save saves a content block to the database
func (r *blockRepository) Save(ctx context.Context, block *blocks.Block) error {
	m := r.ConvertBlockToModel(*block)
	if m.ID == "" {
		uuid := uuid.New()
		m.ID = uuid.String()
		_, err := db.DB.NewInsert().Model(&m).Exec(ctx)
		return err
	}
	_, err := db.DB.NewUpdate().Model(&m).WherePK().Exec(ctx)
	return err
}

// Create saves a content block to the database
func (r *blockRepository) Create(ctx context.Context, block *blocks.Block, locationID string) error {
	contentBlock := models.Block{
		LocationID: locationID,
		Type:       (*block).GetType(),
		Data:       (*block).GetData(),
		Ordering:   1e4,
	}

	uuid := uuid.New()
	contentBlock.ID = uuid.String()
	_, err := db.DB.NewInsert().Model(&contentBlock).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.DB.NewUpdate().Model(&contentBlock).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	*block, err = ConvertModelToBlock(&contentBlock)
	if err != nil {
		return err
	}
	return nil
}

// Update saves a content block to the database
func (r *blockRepository) Update(ctx context.Context, block *blocks.Block) error {
	contentBlock := models.Block{
		ID:         (*block).GetID(),
		Type:       (*block).GetType(),
		Data:       (*block).GetData(),
		LocationID: (*block).GetLocationID(),
		Ordering:   (*block).GetOrder(),
	}
	_, err := db.DB.NewUpdate().Model(&contentBlock).WherePK().Exec(ctx)
	return err
}

// Convert block to model
func (b *blockRepository) ConvertBlockToModel(block blocks.Block) models.Block {
	return models.Block{
		ID:         block.GetID(),
		LocationID: block.GetLocationID(),
		Type:       block.GetType(),
		Ordering:   block.GetOrder(),
		Data:       block.GetData(),
	}
}

func ConvertModelsToBlocks(cbs models.Blocks) (blocks.Blocks, error) {
	b := make(blocks.Blocks, len(cbs))
	for i, cb := range cbs {
		block, err := ConvertModelToBlock(&cb)
		if err != nil {
			return nil, err
		}
		b[i] = block
	}
	return b, nil
}

func ConvertModelToBlock(m *models.Block) (blocks.Block, error) {
	// Convert model to block
	newBlock, err := blocks.CreateFromBaseBlock(blocks.BaseBlock{
		ID:         m.ID,
		LocationID: m.LocationID,
		Type:       m.Type,
		Data:       m.Data,
		Order:      m.Ordering,
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

// Delete deletes a content block from the database
func (r *blockRepository) Delete(ctx context.Context, contentBlockID string) error {
	_, err := db.DB.NewDelete().Model(&models.Block{}).Where("id = ?", contentBlockID).Exec(ctx)
	return err
}

// Reorder reorders the content blocks
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
