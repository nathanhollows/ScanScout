package blocks

import (
	"encoding/json"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
)

func NewBlockFromData(id string, blockType string, data json.RawMessage) (Block, error) {
	switch blockType {
	case "markdown":
		var b MarkdownBlock
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, err
		}
		b.ID = id
		return &b, nil
	default:
		return nil, fmt.Errorf("unknown block type: %s", blockType)
	}
}

func NewBlockFromModel(cb *models.Block) (Block, error) {
	return NewBlockFromData(cb.ID, cb.Type, cb.Data)
}

func NewBlocksFromModel(cbs models.Blocks) (Blocks, error) {
	blocks := make(Blocks, len(cbs))
	for i, cb := range cbs {
		block, err := NewBlockFromModel(&cb)
		if err != nil {
			return nil, err
		}
		blocks[i] = block
	}
	return blocks, nil
}
