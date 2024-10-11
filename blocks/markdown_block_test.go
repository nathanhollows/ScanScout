package blocks

import (
	"encoding/json"
	"testing"

	"github.com/nathanhollows/Rapua/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownBlock_Getters(t *testing.T) {
	block := MarkdownBlock{
		BaseBlock: BaseBlock{
			ID:         "test-id",
			LocationID: "location-123",
			Order:      1,
			Points:     5,
		},
		Content: "Test Content",
	}

	assert.Equal(t, "Markdown", block.GetName())
	assert.Equal(t, "Text written in Markdown.", block.GetDescription())
	assert.Equal(t, "markdown", block.GetType())
	assert.Equal(t, "test-id", block.GetID())
	assert.Equal(t, "location-123", block.GetLocationID())
	assert.Equal(t, 1, block.GetOrder())
	assert.Equal(t, 5, block.GetPoints())
}

func TestMarkdownBlock_ParseData(t *testing.T) {
	data := `{"content":"Test Content"}`
	block := MarkdownBlock{
		BaseBlock: BaseBlock{
			Data: json.RawMessage(data),
		},
	}

	err := block.ParseData()
	require.NoError(t, err)
	assert.Equal(t, "Test Content", block.Content)
}

func TestMarkdownBlock_UpdateBlockData(t *testing.T) {
	block := MarkdownBlock{}
	data := map[string]string{
		"content": "Updated Content",
	}
	err := block.UpdateBlockData(data)
	require.NoError(t, err)
	assert.Equal(t, "Updated Content", block.Content)
}

func TestMarkdownBlock_ValidatePlayerInput(t *testing.T) {
	block := MarkdownBlock{
		BaseBlock: BaseBlock{
			Points: 5,
		},
		Content: "Test Content",
	}

	state := &models.TeamBlockState{
		PlayerData:    nil,
		IsComplete:    false,
		PointsAwarded: 0,
	}

	input := map[string]string{}
	err := block.ValidatePlayerInput(state, input)
	require.NoError(t, err)

	// Assert that state is marked as complete
	assert.True(t, state.IsComplete)
	assert.Equal(t, 0, state.PointsAwarded)
}

func TestMarkdownBlock_CalculatePoints(t *testing.T) {
	block := MarkdownBlock{
		BaseBlock: BaseBlock{
			Points: 5,
		},
		Content: "Test Content",
	}

	input := map[string]string{}
	points, err := block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 0, points) // MarkdownBlock has no points to calculate
}
