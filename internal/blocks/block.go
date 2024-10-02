package blocks

import (
	"encoding/json"
	"fmt"
)

type Block interface {
	GetID() string
	GetType() string
	GetLocationID() string
	GetName() string
	GetDescription() string
	GetOrder() int
	GetIconSVG() string
	GetData() json.RawMessage
	//
	ParseData() error
	UpdateData(data map[string]string) error
	// Render renders the block to html
	Validate(teamCode string, input map[string]string) error
}

type Blocks []Block

type BaseBlock struct {
	ID         string          `json:"-"`
	LocationID string          `json:"-"`
	Type       string          `json:"-"`
	Data       json.RawMessage `json:"-"`
	Order      int             `json:"-"`
}

var RegisteredBlocks = Blocks{
	&MarkdownBlock{},
	&PasswordBlock{},
	// &ChecklistBlock{},
	// &APIBlock{},
}

// CreateFromBaseBlock creates a block from a base block
func CreateFromBaseBlock(baseBlock BaseBlock) (Block, error) {
	switch baseBlock.Type {
	case "markdown":
		return &MarkdownBlock{
			BaseBlock: baseBlock,
		}, nil
	case "password":
		return &PasswordBlock{
			BaseBlock: baseBlock,
		}, nil
	}
	return nil, fmt.Errorf("block type %s not found", baseBlock.Type)
}
