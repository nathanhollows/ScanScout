package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupBlockRepo(t *testing.T) (repositories.BlockRepository, db.Transactor, func()) {
	t.Helper()
	dbc, cleanup := setupDB(t)

	transactor := db.NewTransactor(dbc)

	blockStateRepo := repositories.NewBlockStateRepository(dbc)
	blockRepo := repositories.NewBlockRepository(dbc, blockStateRepo)
	return blockRepo, transactor, cleanup
}

func TestBlockRepository(t *testing.T) {
	repo, transactor, cleanup := setupBlockRepo(t)
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
				tx, err := transactor.BeginTx(context.Background(), &sql.TxOptions{})
				assert.NoError(t, err)
				err = repo.Delete(context.Background(), tx, block.GetID())
				if err != nil {
					tx.Rollback()
					t.Error(err)
				} else {
					tx.Commit()
				}
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
				tx, err := transactor.BeginTx(context.Background(), &sql.TxOptions{})
				assert.NoError(t, err)
				err = repo.Delete(context.Background(), tx, block.GetID())
				if err != nil {
					tx.Rollback()
					t.Error(err)
				} else {
					tx.Commit()
				}
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
				tx, err := transactor.BeginTx(context.Background(), &sql.TxOptions{})
				assert.NoError(t, err)
				err = repo.Delete(context.Background(), tx, block.GetID())
				if err != nil {
					tx.Rollback()
					t.Error(err)
				} else {
					tx.Commit()
				}
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
				tx, err := transactor.BeginTx(context.Background(), &sql.TxOptions{})
				assert.NoError(t, err)
				err = repo.Delete(context.Background(), tx, block.GetID())
				if err != nil {
					tx.Rollback()
					return nil, err
				} else {
					tx.Commit()
					return nil, nil
				}
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

func TestBlockRepository_Bulk(t *testing.T) {
	repo, transactor, cleanup := setupBlockRepo(t)
	defer cleanup()

	tests := []struct {
		name        string
		setup       func() ([]blocks.Block, error)
		action      func(block []blocks.Block) (interface{}, error)
		assertion   func(result interface{}, err error)
		cleanupFunc func(block []blocks.Block)
	}{
		{
			name: "Delete blocks by Location ID",
			setup: func() ([]blocks.Block, error) {
				// Create 3 blocks
				blockSet := make([]blocks.Block, 3)
				locationID := gofakeit.UUID()
				for i := 0; i < 3; i++ {
					block, err := repo.Create(
						context.Background(),
						blocks.NewImageBlock(
							blocks.BaseBlock{
								LocationID: locationID,
								Type:       "image",
								Points:     10,
							},
						),
						locationID,
					)
					if err != nil {
						return nil, err
					}
					blockSet[i] = block
				}
				return blockSet, nil
			},
			action: func(block []blocks.Block) (interface{}, error) {
				tx, err := transactor.BeginTx(context.Background(), &sql.TxOptions{})
				assert.NoError(t, err)

				err = repo.DeleteByLocationID(context.Background(), tx, block[0].GetLocationID())
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				tx.Commit()

				for i, b := range block {
					t.Logf("Checking block %d, ID: %s", i, b.GetID())
					_, err := repo.GetByID(context.Background(), b.GetID())
					if err != nil && err.Error() != "sql: no rows in result set" {
						return nil, err
					}
				}

				return nil, nil
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
			},
			cleanupFunc: func(block []blocks.Block) {},
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
