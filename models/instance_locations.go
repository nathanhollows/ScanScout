package models

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type InstanceLocation struct {
	baseModel

	ID         string `bun:",pk" json:"id"`
	InstanceID string `bun:",notnull" json:"instance_id"`
	LocationID string `bun:",notnull" json:"location_id"`
	CriteriaID string `bun:",notnull" json:"criteria_id"`
	ContentID  string `bun:",notnull" json:"content_id"`

	Instance Instance                `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Location Location                `bun:"rel:has-one,join:location_id=code" json:"location"`
	Criteria CompletionCriteria      `bun:"rel:has-one,join:criteria_id=id" json:"criteria"`
	Content  InstanceLocationContent `bun:"rel:has-one,join:content_id=id" json:"content"`
}

type InstanceLocations []InstanceLocation

// Save saves or updates an instance location
func (i *InstanceLocation) Save(ctx context.Context) error {
	var err error
	if i.ID == "" {
		i.ID = uuid.New().String()
		_, err = db.NewInsert().Model(i).Exec(ctx)
	} else {
		_, err = db.NewUpdate().Model(i).WherePK().Exec(ctx)
	}
	if err != nil {
		log.Error(err)
	}

	return err
}

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

// FindInstanceLocationById finds an instance location by InstanceLocationID
func FindInstanceLocationById(ctx context.Context, id string) (*InstanceLocation, error) {
	var instanceLocation InstanceLocation
	err := db.NewSelect().
		Model(&instanceLocation).
		Where("instance_location.id = ?", id).
		Relation("Location").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &instanceLocation, err
}

// FindInstanceLocationByLocationAndInstance finds an instance location by location and instance
func FindInstanceLocationByLocationAndInstance(ctx context.Context, locationCode, instanceID string) (*InstanceLocation, error) {
	var instanceLocation InstanceLocation
	err := db.NewSelect().
		Model(&instanceLocation).
		Where("instance_location.location_id = ?", locationCode).
		Where("instance_location.instance_id = ?", instanceID).
		Relation("Location").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &instanceLocation, err
}
