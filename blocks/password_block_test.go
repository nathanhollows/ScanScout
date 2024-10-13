package blocks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordBlock_Getters(t *testing.T) {
	block := PasswordBlock{
		BaseBlock: BaseBlock{
			ID:         "test-id",
			LocationID: "location-456",
			Order:      2,
			Points:     10,
		},
		Content:  "Password Content",
		Password: "secret",
		Fuzzy:    true,
	}

	assert.Equal(t, "Password", block.GetName())
	assert.Equal(t, "Players must enter the correct password.", block.GetDescription())
	assert.Equal(t, "password", block.GetType())
	assert.Equal(t, "test-id", block.GetID())
	assert.Equal(t, "location-456", block.GetLocationID())
	assert.Equal(t, 2, block.GetOrder())
	assert.Equal(t, 10, block.GetPoints())
}

func TestPasswordBlock_ParseData(t *testing.T) {
	data := `{"content":"Password Content", "password":"secret", "fuzzy":true}`
	block := PasswordBlock{
		BaseBlock: BaseBlock{
			Data: json.RawMessage(data),
		},
	}

	err := block.ParseData()
	require.NoError(t, err)
	assert.Equal(t, "Password Content", block.Content)
	assert.Equal(t, "secret", block.Password)
	assert.True(t, block.Fuzzy)
}

func TestPasswordBlock_UpdateBlockData(t *testing.T) {
	block := PasswordBlock{}
	data := map[string][]string{
		"content":          {"Updated Password Content"},
		"block-passphrase": {"newsecret"},
		"fuzzy":            {"on"},
	}
	err := block.UpdateBlockData(data)
	require.NoError(t, err)
	assert.Equal(t, "Updated Password Content", block.Content)
	assert.Equal(t, "newsecret", block.Password)
	assert.True(t, block.Fuzzy)
}

func TestPasswordBlock_ValidatePlayerInput(t *testing.T) {
	block := PasswordBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
		Password: "secret",
	}

	state := &mockPlayerState{}

	// Test incorrect password
	input := map[string][]string{"password": {"wrong"}}
	newState, err := block.ValidatePlayerInput(state, input)
	require.NoError(t, err)
	assert.False(t, newState.IsComplete())
	assert.Equal(t, 0, newState.GetPointsAwarded())

	// Test correct password
	input = map[string][]string{"password": {"secret"}}
	newState, err = block.ValidatePlayerInput(state, input)
	require.NoError(t, err)
	assert.True(t, newState.IsComplete())
	assert.Equal(t, 10, newState.GetPointsAwarded())
}

func TestPasswordBlock_CalculatePoints(t *testing.T) {
	block := PasswordBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
	}

	input := map[string][]string{"password": {"any"}}
	points, err := block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 10, points)
}
