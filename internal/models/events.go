package models

type Event struct {
	baseModel

	ID         string `bun:"id,pk,type:varchar(36)"`
	InstanceID string `bun:"instance_id,notnull"`
	Type       string `bun:"type,type:varchar(255)"`
	LocationID string `bun:"location_id,notnull"`
	Points     int    `bun:"points,notnull"`
	Active     bool   `bun:"active,notnull"`

	Instance Instance `bun:"rel:has-one,join:instance_id=id"`
	Location Location `bun:"rel:has-one,join:location_id=id"`
}
