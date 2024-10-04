package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/blocks"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type BlockRepository interface {
	// GetBlocksByLocationID fetches all blocks for a location
	GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error)
	// GetBlockByID fetches a block by its ID
	GetByID(ctx context.Context, blockID string) (blocks.Block, error)
	// Saveblock saves a block to the database
	Save(ctx context.Context, block *blocks.Block) error
	Create(ctx context.Context, block *blocks.Block, locationID string) error
	Update(ctx context.Context, block *blocks.Block) error
	Delete(ctx context.Context, blockID string) error
	Reorder(ctx context.Context, locationID string, blockIDs []string) error
}

type blockRepository struct{}

func NewBlockRepository() BlockRepository {
	return &blockRepository{}
}

// GetByLocationID fetches all blocks for a location
func (r *blockRepository) GetByLocationID(ctx context.Context, locationID string) (blocks.Blocks, error) {
	block := models.Blocks{}
	err := db.DB.NewSelect().
		Model(&block).
		Where("location_id = ?", locationID).
		Order("ordering ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return convertModelsToBlocks(block)
}

// GetByID fetches a block by its ID
func (r *blockRepository) GetByID(ctx context.Context, blockID string) (blocks.Block, error) {
	block := &models.Block{}
	err := db.DB.NewSelect().
		Model(block).
		Where("id = ?", blockID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return convertModelToBlock(block)
}

// Save saves a block to the database
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

// Create saves a block to the database
func (r *blockRepository) Create(ctx context.Context, block *blocks.Block, locationID string) error {
	mBlock := models.Block{
		LocationID:         locationID,
		Type:               (*block).GetType(),
		Data:               (*block).GetData(),
		Ordering:           1e4,
		Points:             (*block).GetPoints(),
		ValidationRequired: (*block).RequiresValidation(),
	}

	uuid := uuid.New()
	mBlock.ID = uuid.String()
	_, err := db.DB.NewInsert().Model(&mBlock).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.DB.NewUpdate().Model(&mBlock).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	*block, err = convertModelToBlock(&mBlock)
	if err != nil {
		return err
	}
	return nil
}

// Update saves a block to the database
func (r *blockRepository) Update(ctx context.Context, block *blocks.Block) error {
	mBlock := models.Block{
		ID:                 (*block).GetID(),
		Type:               (*block).GetType(),
		Data:               (*block).GetData(),
		LocationID:         (*block).GetLocationID(),
		Ordering:           (*block).GetOrder(),
		Points:             (*block).GetPoints(),
		ValidationRequired: (*block).RequiresValidation(),
	}
	_, err := db.DB.NewUpdate().Model(&mBlock).WherePK().Exec(ctx)
	return err
}

// Convert block to model
func (b *blockRepository) ConvertBlockToModel(block blocks.Block) models.Block {
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

func convertModelsToBlocks(cbs models.Blocks) (blocks.Blocks, error) {
	b := make(blocks.Blocks, len(cbs))
	for i, cb := range cbs {
		block, err := convertModelToBlock(&cb)
		if err != nil {
			return nil, err
		}
		b[i] = block
	}
	return b, nil
}

func convertModelToBlock(m *models.Block) (blocks.Block, error) {
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

// Delete deletes a block from the database
func (r *blockRepository) Delete(ctx context.Context, blockID string) error {
	_, err := db.DB.NewDelete().Model(&models.Block{}).Where("id = ?", blockID).Exec(ctx)
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
