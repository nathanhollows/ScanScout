package models

type InstanceLocationContent struct {
	baseModel

	ID                 string           `bun:",pk" json:"id"`
	InstanceLocationID string           `bun:",notnull" json:"instance_location_id"`
	Content            string           `bun:",null" json:"content"`
	InstanceLocation   InstanceLocation `bun:"rel:has-one,join:instance_location_id=id" json:"instance_location"`
}

type InstanceLocationContents []InstanceLocationContent
