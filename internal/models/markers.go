package models

import (
	"context"
	"strings"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/helpers"
)

type Marker struct {
	baseModel

	Code         string  `bun:"code,unique,pk"`
	Lat          float64 `bun:"lat,type:float"`
	Lng          float64 `bun:"lng,type:float"`
	Name         string  `bun:"name,type:varchar(255)"`
	TotalVisits  int     `bun:"total_visits,type:int"`
	CurrentCount int     `bun:"current_count,type:int"`
	AvgDuration  float64 `bun:"avg_duration,type:float"`

	Locations []Location `bun:"rel:has-many,join:code=marker_id"`
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
