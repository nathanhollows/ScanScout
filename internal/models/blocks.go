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

type TeamBlockState struct {
	baseModel
	ID            string          `bun:",pk,notnull" json:"id"`
	TeamCode      string          `bun:",notnull" json:"team_code"`
	BlockID       string          `bun:",notnull" json:"block_id"`
	IsComplete    bool            `bun:",type:bool" json:"is_complete"`
	PointsAwarded int             `bun:",type:int" json:"points"`
	PlayerData    json.RawMessage `bun:",type:jsonb" json:"player_data"`
}
