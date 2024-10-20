package models

import (
	"context"
	"time"

	"github.com/nathanhollows/Rapua/pkg/db"
)

type CheckIn struct {
	baseModel

	InstanceID      string    `bun:"instance_id,notnull"`
	TeamID          string    `bun:"team_code,pk,type:string"`
	LocationID      string    `bun:"location_id,pk,type:string"`
	TimeIn          time.Time `bun:"time_in,type:datetime"`
	TimeOut         time.Time `bun:"time_out,type:datetime"`
	MustCheckOut    bool      `bun:"must_check_out"`
	Points          int       `bun:"points,"`
	BlocksCompleted bool      `bun:"blocks_completed,type:int"`

	Location Location `bun:"rel:has-one,join:location_id=id"`
}

// Save saves or updates a scan
func (s *CheckIn) Save(ctx context.Context) error {
	var err error
	if s.CreatedAt.IsZero() {
		_, err = db.DB.NewInsert().Model(s).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(s).WherePK().Exec(ctx)
	}
	return err
}

// Delete removes the scan from the database
func (s *CheckIn) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(s).WherePK().Exec(ctx)
	return err
}
