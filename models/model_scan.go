package models

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
)

type Scan struct {
	baseModel
	belongsToInstance

	TeamID     string    `bun:",pk,type:string" json:"team_id"`
	LocationID string    `bun:",pk,type:string" json:"location_id"`
	TimeIn     time.Time `bun:",type:datetime" json:"time_in"`
	TimeOut    time.Time `bun:",type:datetime" json:"time_out"`
	Location   Location  `bun:"rel:has-one,join:location_id=code" json:"location"`
}

type Scans []Scan

// Save saves or updates a scan
func (s *Scan) Save() error {
	ctx := context.Background()
	var err error
	if s.CreatedAt.IsZero() {
		_, err = db.NewInsert().Model(s).Exec(ctx)
	} else {
		_, err = db.NewUpdate().Model(s).WherePK().Exec(ctx)
	}
	if err != nil {
		log.Error(err)
	}
	return err
}

// FindScan finds a scan by team and location
func FindScan(teamCode, locationCode string) (*Scan, error) {
	var scan Scan
	err := db.NewSelect().Model(&scan).Where("team_id = ?", teamCode).Where("location_id = ?", locationCode).Scan(context.Background())
	if err != nil {
		log.Error(err)
	}
	return &scan, err
}
