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
	GetPoints() int
	GetIconSVG() string
	GetData() json.RawMessage
	ParseData() error
	UpdateBlockData(data map[string]string) error
	RequiresValidation() bool
	ValidatePlayerInput(input map[string]string) error
	// Calculate partial or full points for a block
	CalculatePoints(input map[string]string) (int, error)
}

type Blocks []Block

type BaseBlock struct {
	ID         string          `json:"-"`
	LocationID string          `json:"-"`
	Type       string          `json:"-"`
	Data       json.RawMessage `json:"-"`
	Order      int             `json:"-"`
	Points     int             `json:"-"`
}

var registeredBlocks = Blocks{
	&MarkdownBlock{},
	&PasswordBlock{},
	// &ChecklistBlock{},
	// &APIBlock{},
}

func GetRegisteredBlocks() Blocks {
	return registeredBlocks
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
