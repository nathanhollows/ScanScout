package models

type Clue struct {
	baseModel

	ID         string `bun:",pk,type:varchar(36)" json:"id"`
	InstanceID string `bun:",notnull" json:"instance_id"`
	LocationID string `bun:",notnull" json:"location_id"`
	Content    string `bun:",type:text" json:"content"`
}
