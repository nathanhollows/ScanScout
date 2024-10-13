package blocks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecklistBlock_ValidatePlayerInput(t *testing.T) {
	// Initial setup for checklist block and mock player state
	block := ChecklistBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
		List: []ChecklistItem{
			{ID: "item-1", Description: "Item 1", Checked: false},
			{ID: "item-2", Description: "Item 2", Checked: false},
		},
	}

	state := &mockPlayerState{}

	// Validate player input where only "item-1" is checked
	input := map[string][]string{
		"checklist-item-ids": {"item-1"}, // "item-1" is checked
	}
	newState, err := block.ValidatePlayerInput(state, input)
	require.NoError(t, err)

	// Assert that player data contains "item-1"
	var playerData checklistPlayerData
	err = json.Unmarshal(state.GetPlayerData(), &playerData)
	require.NoError(t, err)
	assert.Contains(t, playerData.CheckedItems, "item-1")
	assert.False(t, newState.IsComplete())
	assert.Equal(t, 0, newState.GetPointsAwarded())

	// Validate player input where "item-2" is also checked, completing the checklist
	input = map[string][]string{
		"checklist-item-ids": {"item-1", "item-2"}, // "item-1" and "item-2" are checked
	}
	newState, err = block.ValidatePlayerInput(state, input)
	require.NoError(t, err)

	err = json.Unmarshal(newState.GetPlayerData(), &playerData)
	require.NoError(t, err)
	assert.Contains(t, playerData.CheckedItems, "item-1")
	assert.Contains(t, playerData.CheckedItems, "item-2")
	assert.True(t, newState.IsComplete())
	assert.Equal(t, 10, newState.GetPointsAwarded())
}

func TestChecklistBlock_CalculatePoints(t *testing.T) {
	block := ChecklistBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
		List: []ChecklistItem{
			{ID: "item-1", Description: "Item 1", Checked: false},
			{ID: "item-2", Description: "Item 2", Checked: false},
		},
	}

	// Test when not all items are checked
	input := map[string][]string{
		"item-1": {"true"},
	}
	points, err := block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 0, points)

	// Test when all items are checked
	input = map[string][]string{
		"item-1": {"true"},
		"item-2": {"true"},
	}
	points, err = block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 10, points)
}
