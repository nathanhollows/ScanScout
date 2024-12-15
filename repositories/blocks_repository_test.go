package repositories_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupBlockRepo(t *testing.T) (repositories.BlockRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	blockStateRepo := repositories.NewBlockStateRepository(db)
	blockRepo := repositories.NewBlockRepository(db, blockStateRepo)
	return blockRepo, cleanup
}

func TestBlockRepository(t *testing.T) {
	repo, cleanup := setupBlockRepo(t)
	defer cleanup()

	tests := []struct {
		name        string
		setup       func() (blocks.Block, error)
		action      func(block blocks.Block) (interface{}, error)
		assertion   func(result interface{}, err error)
		cleanupFunc func(block blocks.Block)
	}{
		{
			name: "Create new block",
			setup: func() (blocks.Block, error) {
				return repo.Create(
					context.Background(),
					blocks.NewImageBlock(
						blocks.BaseBlock{
							LocationID: gofakeit.UUID(),
							Type:       "image",
							Points:     10,
						},
					),
					gofakeit.UUID(),
				)
			},
			action: func(block blocks.Block) (interface{}, error) {
				return repo.Create(context.Background(), block, gofakeit.UUID())
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
			cleanupFunc: func(block blocks.Block) {
				repo.Delete(context.Background(), block.GetID())
			},
		},
		{
			name: "Get block by ID",
			setup: func() (blocks.Block, error) {
				block, err := repo.Create(
					context.Background(),
					blocks.NewImageBlock(
						blocks.BaseBlock{
							LocationID: gofakeit.UUID(),
							Type:       "image",
							Points:     10,
						},
					),
					gofakeit.UUID(),
				)
				if err != nil {
					return nil, err
				}
				createdBlock, _ := repo.Create(context.Background(), block, gofakeit.UUID())
				return createdBlock.(blocks.Block), nil
			},
			action: func(block blocks.Block) (interface{}, error) {
				return repo.GetByID(context.Background(), block.GetID())
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
			cleanupFunc: func(block blocks.Block) {
				repo.Delete(context.Background(), block.GetID())
			},
		},
		{
			name: "Update block",
			setup: func() (blocks.Block, error) {
				block, err := repo.Create(
					context.Background(),
					blocks.NewImageBlock(
						blocks.BaseBlock{
							LocationID: gofakeit.UUID(),
							Type:       "image",
							Points:     10,
						},
					),
					gofakeit.UUID(),
				)
				if err != nil {
					return nil, err
				}
				createdBlock, _ := repo.Create(context.Background(), block, gofakeit.UUID())
				return createdBlock.(blocks.Block), nil
			},
			action: func(block blocks.Block) (interface{}, error) {
				// This relies on the Image Block
				// TODO: mock the block data
				data := make(map[string][]string)
				data["url"] = []string{"/updated-url"}
				block.UpdateBlockData(data)
				return repo.Update(context.Background(), block)
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
				updatedBlock := result.(blocks.Block)
				assert.Equal(t, `{"content":"/updated-url","caption":"","link":""}`, string(updatedBlock.GetData()))
			},
			cleanupFunc: func(block blocks.Block) {
				repo.Delete(context.Background(), block.GetID())
			},
		},
		{
			name: "Delete block",
			setup: func() (blocks.Block, error) {
				block, err := repo.Create(
					context.Background(),
					blocks.NewImageBlock(
						blocks.BaseBlock{
							LocationID: gofakeit.UUID(),
							Type:       "image",
							Points:     10,
						},
					),
					gofakeit.UUID(),
				)
				if err != nil {
					return nil, err
				}
				createdBlock, _ := repo.Create(context.Background(), block, gofakeit.UUID())
				return createdBlock.(blocks.Block), nil
			},
			action: func(block blocks.Block) (interface{}, error) {
				return nil, repo.Delete(context.Background(), block.GetID())
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
			},
			cleanupFunc: func(block blocks.Block) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := tt.setup()
			assert.NoError(t, err)
			result, err := tt.action(block)
			tt.assertion(result, err)
			if tt.cleanupFunc != nil {
				tt.cleanupFunc(block)
			}
		})
	}
}
