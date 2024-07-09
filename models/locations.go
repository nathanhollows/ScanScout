package models

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type Location struct {
	baseModel

	ID         string `bun:",pk" json:"id"`
	InstanceID string `bun:",notnull" json:"instance_id"`
	CoordsID   string `bun:",notnull" json:"coords_id"`
	CriteriaID string `bun:",notnull" json:"criteria_id"`
	ContentID  string `bun:",notnull" json:"content_id"`

	Instance Instance           `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Coords   Coords             `bun:"rel:has-one,join:coords_id=code" json:"coords"`
	Criteria CompletionCriteria `bun:"rel:has-one,join:criteria_id=id" json:"criteria"`
	Content  LocationContent    `bun:"rel:has-one,join:content_id=id" json:"content"`
}

type Locations []Location

// Save saves or updates an instance location
func (i *Location) Save(ctx context.Context) error {
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
func FindAllInstanceLocations(ctx context.Context) (Locations, error) {
	user := GetUserFromContext(ctx)

	var instanceLocations Locations
	err := db.NewSelect().
		Model(&instanceLocations).
		Where("location.instance_id = ?", user.CurrentInstance.ID).
		Relation("Coords").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return instanceLocations, err
}

// FindInstanceLocationById finds an instance location by InstanceLocationID
func FindInstanceLocationById(ctx context.Context, id string) (*Location, error) {
	var instanceLocation Location
	err := db.NewSelect().
		Model(&instanceLocation).
		Where("location.id = ?", id).
		Relation("Coords").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &instanceLocation, err
}

// FindInstanceLocationByLocationAndInstance finds an instance location by location and instance
func FindInstanceLocationByLocationAndInstance(ctx context.Context, locationCode, instanceID string) (*Location, error) {
	var instanceLocation Location
	err := db.NewSelect().
		Model(&instanceLocation).
		Where("location.coords_id = ?", locationCode).
		Where("location.instance_id = ?", instanceID).
		Relation("Coords").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &instanceLocation, err
}
