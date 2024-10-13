package blocks

import (
	"encoding/json"
	"fmt"
)

type PlayerState interface {
	GetBlockID() string
	GetPlayerID() string
	GetPlayerData() json.RawMessage
	SetPlayerData(data json.RawMessage)
	IsComplete() bool
	SetComplete(complete bool)
	GetPointsAwarded() int
	SetPointsAwarded(points int)
}

type Block interface {
	// Basic Attributes Getters
	GetID() string
	GetType() string
	GetLocationID() string
	GetName() string
	GetDescription() string
	GetOrder() int
	GetPoints() int
	GetIconSVG() string
	GetData() json.RawMessage

	// Data Operations
	ParseData() error
	UpdateBlockData(data map[string][]string) error

	// Validation and Points Calculation
	RequiresValidation() bool
	ValidatePlayerInput(state PlayerState, input map[string][]string) (newState PlayerState, err error)
	CalculatePoints(input map[string][]string) (int, error)
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
	&ChecklistBlock{},
	// &APIBlock{},
}

func GetRegisteredBlocks() Blocks {
	return registeredBlocks
}

func CreateFromBaseBlock(baseBlock BaseBlock) (Block, error) {
	switch baseBlock.Type {
	case "markdown":
		return NewMarkdownBlock(baseBlock), nil
	case "password":
		return NewPasswordBlock(baseBlock), nil
	case "checklist":
		return NewChecklistBlock(baseBlock), nil
	default:
		return nil, fmt.Errorf("block type %s not found", baseBlock.Type)
	}
}

// Example constructor functions
func NewMarkdownBlock(base BaseBlock) *MarkdownBlock {
	return &MarkdownBlock{
		BaseBlock: base,
	}
}

func NewPasswordBlock(base BaseBlock) *PasswordBlock {
	return &PasswordBlock{
		BaseBlock: base,
	}
}

func NewChecklistBlock(base BaseBlock) *ChecklistBlock {
	return &ChecklistBlock{
		BaseBlock: base,
	}
}
