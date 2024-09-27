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

func (b *MarkdownBlock) GetType() string { return "markdown" }
func (b *MarkdownBlock) GetID() string   { return b.ID }
func (b *MarkdownBlock) Data() json.RawMessage {
	data, _ := json.Marshal(b)
	return data
}
