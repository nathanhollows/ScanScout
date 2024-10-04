package models

import (
	"encoding/json"
)

type Block struct {
	ID                 string          `bun:",pk,notnull" json:"id"`
	LocationID         string          `bun:",notnull" json:"location_id"`
	Type               string          `bun:",type:int" json:"type"`
	Data               json.RawMessage `bun:",type:jsonb" json:"data"`
	Ordering           int             `bun:",type:int" json:"order"`
	Points             int             `bun:",type:int" json:"points"`
	ValidationRequired bool            `bun:",type:bool" json:"validation_required"`
}

type Blocks []Block

type TeamBlockState struct {
	baseModel
	TeamCode   string          `bun:",pk,notnull" json:"team_code"`
	BlockID    string          `bun:",pk,notnull" json:"block_id"`
	IsComplete bool            `bun:",type:bool" json:"is_complete"`
	PlayerData json.RawMessage `bun:",type:jsonb" json:"player_data"`
}

type TeamBlockStates []TeamBlockState
