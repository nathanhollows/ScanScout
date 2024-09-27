package blocks

import (
	"context"
	"encoding/json"
	"io"

	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	templates "github.com/nathanhollows/Rapua/internal/templates/blocks"
)

type APIBlock struct {
	ID       string
	EndPoint string // URL that external system will call
	Content  string // Optional content or instructions to display
}

// Ensure APIBlock implements the Block interface
var _ Block = (*APIBlock)(nil)

func (b *APIBlock) Render(ctx context.Context, user models.User, w io.Writer) error {
	html, err := helpers.MarkdownToHTML(b.Content)
	if err != nil {
		return err
	}
	err = templates.MarkdownPlayer(html).Render(ctx, w)
	return err
}

func (b *APIBlock) RenderAdmin(ctx context.Context, user models.User, w io.Writer) error {
	html, err := helpers.MarkdownToHTML(b.Content)
	if err != nil {
		return err
	}
	err = templates.MarkdownAdmin(html).Render(ctx, w)
	return err
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
