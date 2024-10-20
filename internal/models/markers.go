package models

import (
	"context"
	"strings"

	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Marker struct {
	baseModel

	Code         string    `bun:",unique,pk" json:"code"`
	Lat          float64   `bun:",type:float" json:"lat"`
	Lng          float64   `bun:",type:float" json:"lng"`
	Name         string    `bun:",type:varchar(255)" json:"name"`
	TotalVisits  int       `bun:",type:int" json:"total_visits"`
	CurrentCount int       `bun:",type:int" json:"current_count"`
	AvgDuration  float64   `bun:",type:float" json:"avg_duration"`
	Locations    Locations `bun:"rel:has-many,join:code=marker_id" json:"locations"`
}

// Save saves or updates a location
func (l *Marker) Save(ctx context.Context) error {
	insert := false
	var err error
	if l.Code == "" {
		l.Code = helpers.NewCode(5)
		insert = true
	}

	if insert {
		_, err = db.DB.NewInsert().Model(l).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(l).WherePK("code").Exec(ctx)
	}
	return err
}

// FindMarkerByCode returns a marker by code
func FindMarkerByCode(ctx context.Context, code string) (*Marker, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	var marker Marker
	err := db.DB.NewSelect().Model(&marker).Where("code = ?", code).Scan(ctx)
	return &marker, err
}

// SetCoords sets the latitude and longitude of the location
func (l *Marker) SetCoords(lat, lng float64) error {
	l.Lat = lat
	l.Lng = lng
	return nil
}
