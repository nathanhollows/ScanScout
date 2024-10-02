package blocks

import (
	"encoding/json"

	"github.com/nathanhollows/Rapua/internal/models"
)

type APIBlock struct {
	BaseBlock
	// The endpoint to call
	EndPoint string `json:"end_point"`
	// Optional instructions for players
	Content string `json:"content"`
}

func (b *APIBlock) Validate(userID string, input map[string]string) error {
	// No validation required
	return nil
}

func (b *APIBlock) GetName() string { return "API Call" }
func (b *APIBlock) GetDescription() string {
	return "An API call must be made to the specified endpoint."
}
func (b *APIBlock) GetIconSVG() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-terminal"><polyline points="4 17 10 11 4 5"/><line x1="12" x2="20" y1="19" y2="19"/></svg>`
}
func (b *APIBlock) GetType() string { return "api" }
func (b *APIBlock) GetID() string   { return b.ID }
func (b *APIBlock) Data() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

func (b *APIBlock) NewFromBaseBlock(model models.Block) error {
	b.ID = model.ID
	b.LocationID = model.LocationID
	b.Order = model.Order
	err := json.Unmarshal(model.Data, b)
	return err
}
