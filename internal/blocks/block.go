package blocks

import (
	"context"
	"encoding/json"
	"io"

	"github.com/nathanhollows/Rapua/internal/models"
)

type Block interface {
	// Render renders the block to html
	Render(ctx context.Context, user models.User, w io.Writer) error
	RenderAdmin(ctx context.Context, user models.User, w io.Writer) error
	Validate(teamCode string, input map[string]string) error
	GetType() string
	GetID() string
	Data() json.RawMessage
}

type Blocks []Block

var RegisteredBlocks = Blocks{
	&MarkdownBlock{},
	&PasswordBlock{},
	&ChecklistBlock{},
}
