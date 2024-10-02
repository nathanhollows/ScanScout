package models

import "encoding/json"

type Block struct {
	ID         string          `bun:",pk,notnull" json:"id"`
	LocationID string          `bun:",notnull" json:"location_id"`
	Type       string          `bun:",type:int" json:"type"`
	Data       json.RawMessage `bun:",type:jsonb" json:"data"`
	Ordering   int             `bun:",type:int" json:"order"`
}

type Blocks []Block

type TeamBlockProgress struct {
	TeamCode       string          `bun:",pk,notnull" json:"team_code"`
	ContentBlockID string          `bun:",pk,notnull" json:"content_block_id"`
	IsComplete     bool            `bun:",type:bool" json:"is_complete"`
	Progress       json.RawMessage `bun:",type:jsonb" json:"progress"`
}

type TeamBlockProgresses []TeamBlockProgress
