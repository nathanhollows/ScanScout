package blocks

import (
	"encoding/json"
)

type MarkdownBlock struct {
	BaseBlock
	Content string `json:"content"`
}

// Basic Attributes Getters

func (b *MarkdownBlock) GetID() string         { return b.ID }
func (b *MarkdownBlock) GetType() string       { return "markdown" }
func (b *MarkdownBlock) GetLocationID() string { return b.LocationID }
func (b *MarkdownBlock) GetName() string       { return "Markdown" }
func (b *MarkdownBlock) GetDescription() string {
	return "Text written in Markdown."
}
func (b *MarkdownBlock) GetOrder() int  { return b.Order }
func (b *MarkdownBlock) GetPoints() int { return b.Points }
func (b *MarkdownBlock) GetIconSVG() string {
	return `<svg class="w-8 h-8" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M18 8L18 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> <path d="M22 12L18 16L14 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> <path d="M2 16L2 8L6 12L10 8V16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>`
}
func (b *MarkdownBlock) GetData() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}

// Data Operations

func (b *MarkdownBlock) ParseData() error {
	return json.Unmarshal(b.Data, b)
}

func (b *MarkdownBlock) UpdateBlockData(data map[string][]string) error {
	b.Content = data["content"][0]
	return nil
}

// Validation and Points Calculation

func (b *MarkdownBlock) RequiresValidation() bool {
	return false
}

func (b *MarkdownBlock) ValidatePlayerInput(state PlayerState, input map[string][]string) (PlayerState, error) {
	// No validation required for MarkdownBlock; mark as complete
	state.SetComplete(true)
	return state, nil
}

func (b *MarkdownBlock) CalculatePoints(input map[string][]string) (int, error) {
	// MarkdownBlock has no points to calculate
	return 0, nil
}
