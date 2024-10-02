package blocks

import (
	"encoding/json"

	"github.com/nathanhollows/Rapua/internal/models"
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

// Ensure ChecklistBlock implements the Block interface

func (b *ChecklistBlock) Validate(userID string, input map[string]string) error {
	// No validation required
	return nil
}

func (b *ChecklistBlock) GetName() string { return "Checklist" }
func (b *ChecklistBlock) GetDescription() string {
	return "Players must check off all items."
}
func (b *ChecklistBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-list-checks"><path d="m3 17 2 2 4-4"/><path d="m3 7 2 2 4-4"/><path d="M13 6h8"/><path d="M13 12h8"/><path d="M13 18h8"/></svg>`
}
func (b *ChecklistBlock) GetType() string { return "checklist" }
func (b *ChecklistBlock) GetID() string   { return b.ID }
func (b *ChecklistBlock) Data() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

func (b *ChecklistBlock) readFromModel(model models.Block) error {
	b.ID = model.ID
	b.LocationID = model.LocationID
	b.Order = model.Order
	err := json.Unmarshal(model.Data, b)
	return err
}
