package models

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
)

type Location struct {
	baseModel

	ID           string  `bun:",pk,type:integer,autoincrement" json:"id"`
	Code         string  `bun:",notnull" json:"code"`
	Name         string  `bun:",type:varchar(255)" json:"name"`
	InstanceID   string  `bun:",notnull" json:"instance_id"`
	MarkerID     string  `bun:",notnull" json:"marker_id"`
	ContentID    string  `bun:",notnull" json:"content_id"`
	Criteria     string  `bun:",type:varchar(255)" json:"criteria"`
	Order        int     `bun:",type:int" json:"order"`
	TotalVisits  int     `bun:",type:int" json:"total_visits"`
	CurrentCount int     `bun:",type:int" json:"current_count"`
	AvgDuration  float64 `bun:",type:float" json:"avg_duration"`

	Instance Instance        `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Marker   Marker          `bun:"rel:has-one,join:marker_id=code" json:"marker"`
	Content  LocationContent `bun:"rel:has-one,join:content_id=id" json:"content"`
}

type Locations []Location

// Save saves or updates an instance location
func (i *Location) Save(ctx context.Context) error {
	var err error
	if i.Code == "" {
		i.Code = uuid.New().String()
		_, err = db.DB.NewInsert().Model(i).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(i).WherePK().Exec(ctx)
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

	_, err = db.DB.NewDelete().Model(i).WherePK().Exec(ctx)
	return err
}

// FindLocationByInstanceAndCode returns a location by code
func FindLocationByInstanceAndCode(ctx context.Context, instance, code string) (*Location, error) {
	code = strings.ToUpper(code)
	var location Location
	err := db.DB.NewSelect().
		Model(&location).
		Where("marker_id = ?", code).
		Where("instance_id = ?", instance).
		Relation("Marker").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &location, err
}

// FindLocationsByCodes returns a list of locations by code
func FindLocationsByCodes(ctx context.Context, instanceID string, codes []string) Locations {
	var locations Locations
	err := db.DB.NewSelect().
		Model(&locations).
		Where("marker_id in (?)", bun.In(codes)).
		Where("instance_id = ?", instanceID).
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return locations
}

// FindAll returns all locations
func FindAllLocations(ctx context.Context, instanceID string) (Locations, error) {
	var instanceLocations Locations
	err := db.DB.NewSelect().
		Model(&instanceLocations).
		Where("location.instance_id = ?", instanceID).
		Relation("Marker").
		Relation("Content").
		Order("location.order ASC").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return instanceLocations, err
}

// FindInstanceLocationByLocationAndInstance finds an instance location by location and instance
func FindInstanceLocationByLocationAndInstance(ctx context.Context, locationCode, instanceID string) (*Location, error) {
	var instanceLocation Location
	err := db.DB.NewSelect().
		Model(&instanceLocation).
		Where("location.marker_id = ?", locationCode).
		Where("location.instance_id = ?", instanceID).
		Relation("Marker").
		Relation("Content").
		Scan(ctx)
	if err != nil {
		log.Error(err)
	}
	return &instanceLocation, err
}

// FindOrderedLocations returns locations in a specific order
func FindOrderedLocations(ctx context.Context, team *Team) (*Locations, error) {
	// Implement logic to return locations in the order defined by admin
	return nil, nil
}

// FindPseudoRandomLocations returns a set of locations in pseudo-random order
func FindPseudoRandomLocations(ctx context.Context, team *Team) (*Locations, error) {
	// Implement logic to return a set of locations in pseudo-random order
	return nil, nil
}
