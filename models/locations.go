package models

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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

// Delete removes the location from the database
// This will also delete the location content
func (i *Location) Delete(ctx context.Context) error {
	// Delete the location content
	err := i.Content.Delete(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewDelete().Model(i).WherePK().Exec(ctx)
	return err
}

// FindLocationByCode returns a location by code
func FindLocationByCode(ctx context.Context, code string) (*Location, error) {
	code = strings.ToUpper(code)
	var location Location
	err := db.NewSelect().
		Model(&location).
		Where("coords_id = ?", code).
		Where("instance_id = ?", GetUserFromContext(ctx).CurrentInstanceID).
		Relation("Coords").
		Relation("Criteria").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &location, err
}

// FindLocationsByCodes returns a list of locations by code
func FindLocationsByCodes(ctx context.Context, codes []string) Locations {
	var locations Locations
	err := db.NewSelect().
		Model(&locations).
		Where("coords_id in (?)", bun.In(codes)).
		Where("instance_id = ?", GetUserFromContext(ctx).CurrentInstance.ID).
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return locations
}

// FindAll returns all locations
func FindAllLocations(ctx context.Context) (Locations, error) {
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
