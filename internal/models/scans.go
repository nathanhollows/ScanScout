package models

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Scan struct {
	baseModel

	InstanceID  string    `bun:",notnull" json:"instance_id"`
	TeamID      string    `bun:",pk,type:string" json:"team_id"`
	LocationID  string    `bun:",pk,type:string" json:"location_id"`
	TimeIn      time.Time `bun:",type:datetime" json:"time_in"`
	TimeOut     time.Time `bun:",type:datetime" json:"time_out"`
	Location    Location  `bun:"rel:has-one,join:location_id=code" json:"location"`
	MustScanOut int       `bun:"" json:"must_scan_out"`
}

type Scans []Scan

// Save saves or updates a scan
func (s *Scan) Save(ctx context.Context) error {
	var err error
	if s.CreatedAt.IsZero() {
		_, err = db.DB.NewInsert().Model(s).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(s).WherePK().Exec(ctx)
	}
	if err != nil {
		log.Error(err)
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
	if err != nil {
		log.Error(err)
	}
	return &scan, err
}
