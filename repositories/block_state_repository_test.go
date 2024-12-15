package repositories_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nathanhollows/Rapua/blocks"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/stretchr/testify/assert"
)

func setupBlockStateRepo(t *testing.T) (repositories.BlockStateRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	blockStateRepository := repositories.NewBlockStateRepository(db)
	return blockStateRepository, cleanup
}

func TestBlockStateRepository(t *testing.T) {
	repo, cleanup := setupBlockStateRepo(t)
	defer cleanup()

	tests := []struct {
		name        string
		setup       func() (blocks.PlayerState, error)
		action      func(state blocks.PlayerState) (interface{}, error)
		assertion   func(result interface{}, err error)
		cleanupFunc func(state blocks.PlayerState)
	}{
		{
			name: "Create new player state",
			setup: func() (blocks.PlayerState, error) {
				return repo.NewBlockState(context.Background(), gofakeit.UUID(), gofakeit.UUID())
			},
			action: func(state blocks.PlayerState) (interface{}, error) {
				return repo.Create(context.Background(), state)
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
			cleanupFunc: func(state blocks.PlayerState) {
				repo.Delete(context.Background(), state.GetBlockID(), state.GetPlayerID())
			},
		},
		{
			name: "Get player state by block and team",
			setup: func() (blocks.PlayerState, error) {
				state, _ := repo.NewBlockState(context.Background(), gofakeit.UUID(), gofakeit.UUID())
				return repo.Create(context.Background(), state)
			},
			action: func(state blocks.PlayerState) (interface{}, error) {
				return repo.GetByBlockAndTeam(context.Background(), state.GetBlockID(), state.GetPlayerID())
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
			cleanupFunc: func(state blocks.PlayerState) {
				repo.Delete(context.Background(), state.GetBlockID(), state.GetPlayerID())
			},
		},
		{
			name: "Update player state",
			setup: func() (blocks.PlayerState, error) {
				state, _ := repo.NewBlockState(context.Background(), gofakeit.UUID(), gofakeit.UUID())
				createdState, _ := repo.Create(context.Background(), state)
				return createdState, nil
			},
			action: func(state blocks.PlayerState) (interface{}, error) {
				state.SetPlayerData([]byte(`{"key":"value"}`))
				state.SetComplete(true)
				state.SetPointsAwarded(100)
				return repo.Update(context.Background(), state)
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
				updatedState := result.(blocks.PlayerState)
				assert.Equal(t, true, updatedState.IsComplete())
				assert.Equal(t, 100, updatedState.GetPointsAwarded())
			},
			cleanupFunc: func(state blocks.PlayerState) {
				repo.Delete(context.Background(), state.GetBlockID(), state.GetPlayerID())
			},
		},
		{
			name: "Delete player state",
			setup: func() (blocks.PlayerState, error) {
				state, _ := repo.NewBlockState(context.Background(), gofakeit.UUID(), gofakeit.UUID())
				return repo.Create(context.Background(), state)
			},
			action: func(state blocks.PlayerState) (interface{}, error) {
				return nil, repo.Delete(context.Background(), state.GetBlockID(), state.GetPlayerID())
			},
			assertion: func(result interface{}, err error) {
				assert.NoError(t, err)
			},
			cleanupFunc: func(state blocks.PlayerState) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := tt.setup()
			assert.NoError(t, err)
			result, err := tt.action(state)
			tt.assertion(result, err)
			if tt.cleanupFunc != nil {
				tt.cleanupFunc(state)
			}
		})
	}
}
