package blocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYoutubeBlock_Getters(t *testing.T) {
	block := YoutubeBlock{
		BaseBlock: BaseBlock{
			ID:         "test-id",
			LocationID: "location-123",
			Order:      1,
			Points:     5,
		},
		URL: "https://www.youtube.com/watch?v=12345",
	}

	assert.Equal(t, "Youtube", block.GetName())
	assert.Equal(t, "youtube", block.GetType())
	assert.Equal(t, "test-id", block.GetID())
	assert.Equal(t, "location-123", block.GetLocationID())
	assert.Equal(t, 1, block.GetOrder())
	assert.Equal(t, 5, block.GetPoints())
}

func TestYoutubeBlock_UpdateBlockData(t *testing.T) {
	block := YoutubeBlock{}
	data := map[string][]string{
		"URL": {"https://www.youtube.com/watch?v=54321"},
	}
	err := block.UpdateBlockData(data)
	require.NoError(t, err)
	assert.Equal(t, "https://www.youtube.com/watch?v=54321", block.URL)
}

func TestYoutubeBlock_ValidatePlayerInput(t *testing.T) {
	block := YoutubeBlock{
		BaseBlock: BaseBlock{
			Points: 5,
		},
		URL: "https://www.youtube.com/watch?v=12345",
	}

	state := &mockPlayerState{}

	input := map[string][]string{}
	newState, err := block.ValidatePlayerInput(state, input)
	require.NoError(t, err)

	// Assert that state is marked as complete
	assert.True(t, newState.IsComplete())
	assert.Equal(t, 0, newState.GetPointsAwarded())
}

func TestYoutubeBlock_CalculatePoints(t *testing.T) {
	block := YoutubeBlock{
		BaseBlock: BaseBlock{
			Points: 5,
		},
		URL: "https://www.youtube.com/watch?v=12345",
	}

	input := map[string][]string{}
	points, err := block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 0, points) // YoutubeBlock has no points to calculate
}
