package blocks

import (
	"context"
	"encoding/json"
	"io"

	"github.com/nathanhollows/Rapua/internal/models"
)

type Block interface {
	GetID() string
	GetType() string
	GetName() string
	GetDescription() string
	GetIconSVG() string
	// Data returns the block data as a json.RawMessage
	Data() json.RawMessage
	// Render renders the block to html
	Render(ctx context.Context, user models.User, w io.Writer) error
	RenderAdmin(ctx context.Context, user models.User, w io.Writer) error
	Validate(teamCode string, input map[string]string) error
}

type Blocks []Block

var RegisteredBlocks = Blocks{
	&MarkdownBlock{},
	&PasswordBlock{},
	&ChecklistBlock{},
	&APIBlock{},
}
