package blocks_test

import (
	"encoding/json"
	"testing"

	"github.com/nathanhollows/Rapua/blocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateFromBaseBlock(t *testing.T) {
	t.Run("creates MarkdownBlock from base block", func(t *testing.T) {
		baseBlock := blocks.BaseBlock{
			ID:         "1",
			LocationID: "loc1",
			Type:       "markdown",
			Data:       json.RawMessage(`{"content": "Hello World"}`),
			Order:      1,
			Points:     10,
		}

		block, err := blocks.CreateFromBaseBlock(baseBlock)
		assert.NoError(t, err)
		assert.IsType(t, &blocks.MarkdownBlock{}, block)
		assert.Equal(t, "markdown", block.GetType())
		assert.Equal(t, "1", block.GetID())
		assert.Equal(t, 10, block.GetPoints())
	})

	t.Run("creates AnswerBlock from base block", func(t *testing.T) {
		baseBlock := blocks.BaseBlock{
			ID:         "2",
			LocationID: "loc2",
			Type:       "answer",
			Data:       json.RawMessage(`{"answer": "secret"}`),
			Order:      2,
			Points:     20,
		}

		block, err := blocks.CreateFromBaseBlock(baseBlock)
		assert.NoError(t, err)
		assert.IsType(t, &blocks.AnswerBlock{}, block)
		assert.Equal(t, "answer", block.GetType())
		assert.Equal(t, "2", block.GetID())
		assert.Equal(t, 20, block.GetPoints())
	})

	t.Run("returns error for unregistered block type", func(t *testing.T) {
		baseBlock := blocks.BaseBlock{
			ID:         "3",
			LocationID: "loc3",
			Type:       "unknown",
			Data:       json.RawMessage(`{}`),
			Order:      3,
			Points:     30,
		}

		block, err := blocks.CreateFromBaseBlock(baseBlock)
		assert.Error(t, err)
		assert.Nil(t, block)
		assert.EqualError(t, err, "block type unknown not found")
	})
}

func TestGetRegisteredBlocks(t *testing.T) {
	t.Run("returns registered blocks", func(t *testing.T) {
		blocklist := blocks.GetRegisteredBlocks()
		assert.Len(t, blocklist, 6)
		assert.IsType(t, &blocks.MarkdownBlock{}, blocklist[0])
		assert.IsType(t, &blocks.ImageBlock{}, blocklist[1])
		assert.IsType(t, &blocks.AnswerBlock{}, blocklist[2])
		assert.IsType(t, &blocks.PincodeBlock{}, blocklist[3])
		assert.IsType(t, &blocks.ChecklistBlock{}, blocklist[4])
		assert.IsType(t, &blocks.YoutubeBlock{}, blocklist[5])
	})
}
