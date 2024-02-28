package models

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
)

type Scan struct {
	baseModel
	belongsToInstance

	TeamID     string    `bun:",type:string" json:"team_id"`
	LocationID string    `bun:",type:string" json:"location_id"`
	TimeIn     time.Time `bun:",type:datetime" json:"time_in"`
	TimeOut    time.Time `bun:",type:datetime" json:"time_out"`
}

type Scans []Scan

// Save saves or updates a scan
func (s *Scan) Save() error {
	ctx := context.Background()
	_, err := db.NewInsert().Model(s).Exec(ctx)
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
