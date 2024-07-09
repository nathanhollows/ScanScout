package models

type InstanceLocationContent struct {
	baseModel

	ID                 string `bun:",pk" json:"id"`
	InstanceLocationID string `bun:",notnull" json:"instance_location_id"`
	Content            string `bun:",null" json:"content"`
}

type InstanceLocationContents []InstanceLocationContent
