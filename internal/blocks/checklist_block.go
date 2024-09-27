package blocks

import (
	"context"
	"encoding/json"
	"io"

	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	templates "github.com/nathanhollows/Rapua/internal/templates/blocks"
)

type ChecklistBlock struct {
	ID      string
	Content string
}

// Ensure ChecklistBlock implements the Block interface
var _ Block = (*ChecklistBlock)(nil)

func (b *ChecklistBlock) Render(ctx context.Context, user models.User, w io.Writer) error {
	html, err := helpers.MarkdownToHTML(b.Content)
	if err != nil {
		return err
	}
	err = templates.MarkdownPlayer(html).Render(ctx, w)
	return err
}

func (b *ChecklistBlock) RenderAdmin(ctx context.Context, user models.User, w io.Writer) error {
	html, err := helpers.MarkdownToHTML(b.Content)
	if err != nil {
		return err
	}
	err = templates.MarkdownAdmin(html).Render(ctx, w)
	return err
}

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
