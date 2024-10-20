package models

import (
	"context"
	"strings"
	"time"

	"github.com/nathanhollows/Rapua/pkg/db"
)

type Scan struct {
	baseModel

	InstanceID      string    `bun:"instance_id,notnull"`
	TeamID          string    `bun:"team_id,pk,type:string"`
	LocationID      string    `bun:"location_id,pk,type:string"`
	TimeIn          time.Time `bun:"time_in,type:datetime"`
	TimeOut         time.Time `bun:"time_out,type:datetime"`
	MustScanOut     bool      `bun:"must_scan_out"`
	Points          int       `bun:"points,"`
	BlocksCompleted bool      `bun:"blocks_completed,type:int"`

	Location Location `bun:"rel:has-one,join:location_id=id"`
}

// Save saves or updates a scan
func (s *Scan) Save(ctx context.Context) error {
	var err error
	if s.CreatedAt.IsZero() {
		_, err = db.DB.NewInsert().Model(s).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(s).WherePK().Exec(ctx)
	}
	return err
}

// Delete removes the scan from the database
func (s *Scan) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(s).WherePK().Exec(ctx)
	return err
}

// FindScan finds a scan by team and location
func FindScan(ctx context.Context, teamCode, locationCode string) (*Scan, error) {
	teamCode = strings.ToUpper(teamCode)
	locationCode = strings.ToUpper(locationCode)
	var scan Scan
	err := db.DB.NewSelect().Model(&scan).Where("team_id = ?", teamCode).Where("location_id = ?", locationCode).Scan(ctx)
	return &scan, err
}

// String returns a string representation of a scan
func (s *Scan) String() string {
	return s.TeamID + " " + s.LocationID
}
