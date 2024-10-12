package blocks

import (
	"encoding/json"
	"fmt"

	"github.com/nathanhollows/Rapua/models"
)

type ChecklistBlock struct {
	BaseBlock
	Content string          `json:"content"`
	List    []ChecklistItem `json:"list"`
}

type ChecklistItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Checked     bool   `json:"checked"`
}

// Unexported struct for storing player progress data in a block
type checklistPlayerData struct {
	CheckedItems []string `json:"checked_items"`
}

// Ensure ChecklistBlock implements the Block interface

// Basic Attributes Getters
func (b *ChecklistBlock) GetName() string { return "Checklist" }

func (b *ChecklistBlock) GetDescription() string {
	return "Players must check off all items."
}

func (b *ChecklistBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-list-checks"><path d="m3 17 2 2 4-4"/><path d="m3 7 2 2 4-4"/><path d="M13 6h8"/><path d="M13 12h8"/><path d="M13 18h8"/></svg>`
}

func (b *ChecklistBlock) GetType() string { return "checklist" }

func (b *ChecklistBlock) GetID() string { return b.ID }

func (b *ChecklistBlock) GetLocationID() string { return b.LocationID }

func (b *ChecklistBlock) GetOrder() int { return b.Order }

func (b *ChecklistBlock) GetPoints() int { return b.Points }

func (b *ChecklistBlock) GetData() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

// Data Operations
func (b *ChecklistBlock) ParseData() error {
	return json.Unmarshal(b.Data, b)
}

	b.Content = data["content"]
func (b *ChecklistBlock) UpdateBlockData(data map[string][]string) error {
	return nil
}

// Validation and Points Calculation
func (b *ChecklistBlock) RequiresValidation() bool { return true }

func (b *ChecklistBlock) ValidatePlayerInput(state *models.TeamBlockState, input map[string]string) error {
	// Parse player data from the existing state
	var playerData checklistPlayerData
	if state.PlayerData != nil {
		err := json.Unmarshal(state.PlayerData, &playerData)
		if err != nil {
			return fmt.Errorf("failed to parse player data: %w", err)
		}
	}

	// Update checked items based on the player's input
	for itemID, checked := range input {
		if checked == "true" {
			// Only add unique items to CheckedItems
			itemAlreadyChecked := false
			for _, existingID := range playerData.CheckedItems {
				if existingID == itemID {
					itemAlreadyChecked = true
					break
				}
			}
			if !itemAlreadyChecked {
				playerData.CheckedItems = append(playerData.CheckedItems, itemID)
			}
		}
	}

	// Marshal the updated player data back into the state
	newPlayerData, err := json.Marshal(playerData)
	if err != nil {
		return fmt.Errorf("failed to save player data: %w", err)
	}
	state.PlayerData = newPlayerData

	// Mark the state as complete if all items are checked
	allChecked := true
	for _, item := range b.List {
		itemChecked := false
		for _, checkedID := range playerData.CheckedItems {
			if item.ID == checkedID {
				itemChecked = true
				break
			}
		}
		if !itemChecked {
			allChecked = false
			break
		}
	}

	if allChecked {
		state.IsComplete = true
		state.PointsAwarded = b.Points
	}

	return nil
}

func (b *ChecklistBlock) CalculatePoints(input map[string]string) (int, error) {
	// For ChecklistBlock, return full points if all items are checked, otherwise 0 points
	allChecked := true
	for _, item := range b.List {
		if input[item.ID] != "true" {
			allChecked = false
			break
		}
	}
	if allChecked {
		return b.Points, nil
	}
	return 0, nil
}

// Additional methods for reading from models
func (b *ChecklistBlock) readFromModel(model models.Block) error {
	b.ID = model.ID
	b.LocationID = model.LocationID
	b.Order = model.Ordering
	return json.Unmarshal(model.Data, b)
}
