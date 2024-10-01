package blocks

import (
	"encoding/json"

	"github.com/nathanhollows/Rapua/internal/models"
)

type Block interface {
	GetID() string
	GetType() string
	GetLocationID() string
	GetName() string
	GetDescription() string
	GetIconSVG() string
	GetAdminData() interface{}
	// Data returns the block data as a json.RawMessage
	Data() json.RawMessage
	// Render renders the block to html
	Validate(teamCode string, input map[string]string) error
	readFromModel(model models.Block) error
}

type Blocks []Block

type BaseBlock struct {
	ID         string `json:"-"`
	LocationID string `json:"-"`
	Type       string `json:"-"`
	Order      int    `json:"-"`
}

var RegisteredBlocks = Blocks{
	&MarkdownBlock{},
	&PasswordBlock{},
	// &ChecklistBlock{},
	// &APIBlock{},
}
