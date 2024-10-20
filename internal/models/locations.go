package models

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Location struct {
	baseModel

	ID           string           `bun:"id,pk,notnull"`
	Name         string           `bun:"name,type:varchar(255)"`
	InstanceID   string           `bun:"instance_id,notnull"`
	MarkerID     string           `bun:"marker_id,notnull"`
	ContentID    string           `bun:"content_id,notnull"`
	Criteria     string           `bun:"criteria,type:varchar(255)"`
	Order        int              `bun:"order,type:int"`
	TotalVisits  int              `bun:"total_visits,type:int"`
	CurrentCount int              `bun:"current_count,type:int"`
	AvgDuration  float64          `bun:"avg_duration,type:float"`
	Completion   CompletionMethod `bun:"completion,type:int"`
	Points       int              `bun:"points,"`

	Clues    []models.Clue  `bun:"rel:has-many,join:id=location_id"`
	Instance Instance       `bun:"rel:has-one,join:instance_id=id"`
	Marker   Marker         `bun:"rel:has-one,join:marker_id=code"`
	Blocks   []models.Block `bun:"rel:has-many,join:id=location_id"`
}

type Locations []Location

// Save saves or updates an instance location
func (i *Location) Save(ctx context.Context) error {
	var err error
	if i.ID == "" {
		i.ID = uuid.New().String()
		_, err = db.DB.NewInsert().Model(i).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(i).WherePK().Exec(ctx)
	}

	return err
}

// Delete removes the location from the database
func (i *Location) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(i).WherePK().Exec(ctx)
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
		Scan(ctx)
	return &location, err
}

// FindAll returns all locations
func FindAllLocations(ctx context.Context, instanceID string) (Locations, error) {
	var instanceLocations Locations
	err := db.DB.NewSelect().
		Model(&instanceLocations).
		Where("location.instance_id = ?", instanceID).
		Relation("Marker").
		Order("location.order ASC").
		Scan(ctx)
	return instanceLocations, err
}
