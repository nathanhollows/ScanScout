package models

type Clue struct {
	baseModel

	ID         string `bun:"id,pk,type:varchar(36)"`
	InstanceID string `bun:"instance_id,notnull"`
	LocationID string `bun:"location_id,notnull"`
	Content    string `bun:"content,type:text"`
}
