package models

import (
	"context"

	"github.com/charmbracelet/log"
)

type InstanceLocation struct {
	baseModel

	ID         string                  `bun:",pk" json:"id"`
	InstanceID string                  `bun:",notnull" json:"instance_id"`
	LocationID string                  `bun:",notnull" json:"location_id"`
	CriteriaID string                  `bun:",notnull" json:"criteria_id"`
	Instance   Instance                `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Location   Location                `bun:"rel:has-one,join:location_id=code" json:"location"`
	Criteria   CompletionCriteria      `bun:"rel:has-one,join:criteria_id=id" json:"criteria"`
	Content    InstanceLocationContent `bun:"rel:has-one,join:id=instance_location_id" json:"instance_location_content"`
}

type InstanceLocations []InstanceLocation

// FindAll returns all locations
func FindAllInstanceLocations(ctx context.Context) (InstanceLocations, error) {
	user := GetUserFromContext(ctx)

	var instanceLocations InstanceLocations
	err := db.NewSelect().
		Model(&instanceLocations).
		Where("instance_location.instance_id = ?", user.CurrentInstance.ID).
		Relation("Location").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return instanceLocations, err
}
