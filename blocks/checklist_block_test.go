package blocks

import (
	"encoding/json"
	"testing"

	"github.com/nathanhollows/Rapua/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecklistBlock_Getters(t *testing.T) {
	block := ChecklistBlock{
		BaseBlock: BaseBlock{
			ID:         "test-id",
			LocationID: "location-123",
			Order:      1,
			Points:     10,
		},
		Content: "Test Content",
	}

	assert.Equal(t, "Checklist", block.GetName())
	assert.Equal(t, "Players must check off all items.", block.GetDescription())
	assert.Equal(t, "checklist", block.GetType())
	assert.Equal(t, "test-id", block.GetID())
	assert.Equal(t, "location-123", block.GetLocationID())
	assert.Equal(t, 1, block.GetOrder())
	assert.Equal(t, 10, block.GetPoints())
}

func TestChecklistBlock_ParseData(t *testing.T) {
	data := `{"content":"Test Content", "list":[{"id":"item-1", "description":"Item 1", "checked":false}]}`
	block := ChecklistBlock{
		BaseBlock: BaseBlock{
			Data: json.RawMessage(data),
		},
	}

	err := block.ParseData()
	require.NoError(t, err)
	assert.Equal(t, "Test Content", block.Content)
	assert.Len(t, block.List, 1)
	assert.Equal(t, "item-1", block.List[0].ID)
	assert.Equal(t, "Item 1", block.List[0].Description)
	assert.False(t, block.List[0].Checked)
}

func TestChecklistBlock_UpdateBlockData(t *testing.T) {
	block := ChecklistBlock{}
	data := map[string][]string{
		"content":            {"Updated Content"},
		"checklist-items":    {"Get Milk", "Get Eggs", "Get Bread"},
		"checklist-item-ids": {"", "", "existing-id"},
	}
	err := block.UpdateBlockData(data)
	require.NoError(t, err)
	assert.Equal(t, "Updated Content", block.Content)
	assert.Len(t, block.List, 3)
	assert.NotEmpty(t, block.List[0].ID) // ID should be generated
	assert.Equal(t, "Get Milk", block.List[0].Description)
	assert.NotEmpty(t, block.List[1].ID) // ID should be generated
	assert.Equal(t, "Get Eggs", block.List[1].Description)
	assert.Equal(t, "existing-id", block.List[2].ID) // Preserved ID
	assert.Equal(t, "Get Bread", block.List[2].Description)
}

func TestChecklistBlock_ValidatePlayerInput(t *testing.T) {
	// Initial setup for checklist block and team state
	block := ChecklistBlock{
		BaseBlock: BaseBlock{
			Points: 10,
		},
		List: []ChecklistItem{
			{ID: "item-1", Description: "Item 1", Checked: false},
			{ID: "item-2", Description: "Item 2", Checked: false},
		},
	}

	state := &models.TeamBlockState{
		PlayerData:    nil,
		IsComplete:    false,
		PointsAwarded: 0,
	}

	// Validate player input where only "item-1" is checked
	input := map[string][]string{
		"checklist-item-ids": {"item-1"}, // "item-1" is checked
	}
	err := block.ValidatePlayerInput(state, input)
	require.NoError(t, err)

	// Assert that player data contains "item-1"
	var playerData checklistPlayerData
	err = json.Unmarshal(state.PlayerData, &playerData)
	require.NoError(t, err)
	assert.Contains(t, playerData.CheckedItems, "item-1")
	assert.False(t, state.IsComplete)
	assert.Equal(t, 0, state.PointsAwarded)

	// Validate player input where "item-2" is also checked, completing the checklist
	input = map[string][]string{
		"checklist-item-ids": {"item-1", "item-2"}, // "item-1" and "item-2" are checked
	}
	err = block.ValidatePlayerInput(state, input)
	require.NoError(t, err)

	err = json.Unmarshal(state.PlayerData, &playerData)
	require.NoError(t, err)
	assert.Contains(t, playerData.CheckedItems, "item-1")
	assert.Contains(t, playerData.CheckedItems, "item-2")
	assert.True(t, state.IsComplete)
	assert.Equal(t, 10, state.PointsAwarded)
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
	input := map[string]string{
		"item-1": "true",
	}
	points, err := block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 0, points)

	// Test when all items are checked
	input = map[string]string{
		"item-1": "true",
		"item-2": "true",
	}
	points, err = block.CalculatePoints(input)
	require.NoError(t, err)
	assert.Equal(t, 10, points)
}
