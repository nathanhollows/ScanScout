package blocks

import (
	"context"
	"encoding/json"
	"io"

	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	templates "github.com/nathanhollows/Rapua/internal/templates/blocks"
)

type MarkdownBlock struct {
	ID      string
	Content string
}

// Ensure MarkdownBlock implements the Block interface
var _ Block = (*MarkdownBlock)(nil)

func (b *MarkdownBlock) Render(ctx context.Context, user models.User, w io.Writer) error {
	html, err := helpers.MarkdownToHTML(b.Content)
	if err != nil {
		return err
	}
	err = templates.MarkdownPlayer(html).Render(ctx, w)
	return err
}

func (b *MarkdownBlock) RenderAdmin(ctx context.Context, user models.User, w io.Writer) error {
	html, err := helpers.MarkdownToHTML(b.Content)
	if err != nil {
		return err
	}
	err = templates.MarkdownAdmin(html).Render(ctx, w)
	return err
}

func (b *MarkdownBlock) Validate(userID string, input map[string]string) error {
	// No validation required
	return nil
}

func (b *MarkdownBlock) GetName() string { return "Markdown" }
func (b *MarkdownBlock) GetDescription() string {
	return "Text written in Markdown."
}
func (b *MarkdownBlock) GetIconSVG() string {
	return `<svg class="w-8 h-8" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M18 8L18 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> <path d="M22 12L18 16L14 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> <path d="M2 16L2 8L6 12L10 8V16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg> `
}
func (b *MarkdownBlock) GetType() string { return "markdown" }
func (b *MarkdownBlock) GetID() string   { return b.ID }
func (b *MarkdownBlock) Data() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}
