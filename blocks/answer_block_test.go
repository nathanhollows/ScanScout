package blocks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnswerBlock_Getters(t *testing.T) {
	block := AnswerBlock{
		BaseBlock: BaseBlock{
			ID:         "test-id",
			LocationID: "location-456",
			Order:      2,
			Points:     10,
		},
		Prompt: "Answer Content",
		Answer: "secret",
		Fuzzy:  true,
	}

	assert.Equal(t, "Answer", block.GetName())
	assert.Equal(t, "answer", block.GetType())
	assert.Equal(t, "test-id", block.GetID())
	assert.Equal(t, "location-456", block.GetLocationID())
	assert.Equal(t, 2, block.GetOrder())
	assert.Equal(t, 10, block.GetPoints())
}

func TestAnswerBlock_ParseData(t *testing.T) {
	data := `{"prompt":"Answer Content", "answer":"secret", "fuzzy":true}`
	block := AnswerBlock{
		BaseBlock: BaseBlock{
			Data: json.RawMessage(data),
		},
	}

	err := block.ParseData()
	require.NoError(t, err)
	assert.Equal(t, "Answer Content", block.Prompt)
	assert.Equal(t, "secret", block.Answer)
	assert.True(t, block.Fuzzy)
}

func TestAnswerBlock_UpdateBlockData(t *testing.T) {
	// Update all fields
	block := AnswerBlock{}
	data := map[string][]string{
		"prompt": {"Updated Answer Content"},
		"answer": {"newsecret"},
		"fuzzy":  {"on"},
	}
	err := block.UpdateBlockData(data)
	require.NoError(t, err)
	assert.Equal(t, "Updated Answer Content", block.Prompt)
	assert.Equal(t, "newsecret", block.Answer)
	assert.True(t, block.Fuzzy)

	// Update without fuzzy
	block = AnswerBlock{}
	data = map[string][]string{
		"prompt": {"Updated Answer Content"},
		"answer": {"newsecret"},
	}
	err = block.UpdateBlockData(data)
	require.NoError(t, err)
	assert.Equal(t, "Updated Answer Content", block.Prompt)
	assert.Equal(t, "newsecret", block.Answer)
	assert.False(t, block.Fuzzy)
}

func TestAnswerBlock_ValidatePlayerInput(t *testing.T) {
	block := AnswerBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
		Answer: "secret",
	}

	state := &mockPlayerState{}

	// Test incorrect answer
	input := map[string][]string{"answer": {"wrong"}}
	newState, err := block.ValidatePlayerInput(state, input)
	require.NoError(t, err)
	assert.False(t, newState.IsComplete())
	assert.Equal(t, 0, newState.GetPointsAwarded())

	// Test correct answer
	input = map[string][]string{"answer": {"secret"}}
	newState, err = block.ValidatePlayerInput(state, input)
	require.NoError(t, err)
	assert.True(t, newState.IsComplete())
	assert.Equal(t, 10, newState.GetPointsAwarded())
}

func TestAnswerBlock_CalculatePoints(t *testing.T) {
	block := AnswerBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
	}

	input := map[string][]string{"answer": {"any"}}
	points, err := block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 10, points)
}
