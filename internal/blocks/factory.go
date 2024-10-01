package blocks

import (
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
)

func ConvertModelsToBlocks(cbs models.Blocks) (Blocks, error) {
	blocks := make(Blocks, len(cbs))
	for i, cb := range cbs {
		block, err := ConvertModelToBlock(&cb)
		if err != nil {
			return nil, err
		}
		blocks[i] = block
	}
	return blocks, nil
}

func ConvertModelToBlock(m *models.Block) (Block, error) {
	var newBlock Block
	for _, rb := range RegisteredBlocks {
		if rb.GetType() == m.Type {
			newBlock = rb
			break
		}
	}
	if newBlock == nil {
		return nil, fmt.Errorf("unknown block type: %s", m.Type)
	}
	err := newBlock.readFromModel(*m)
	if err != nil {
		return nil, err
	}
	return newBlock, nil
}
