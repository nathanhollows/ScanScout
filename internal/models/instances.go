package models

import (
	"time"

	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

// Instance represents a single planned activity belonging to a user
// Instance is used to match users, teams, locations, and scans
type Instance struct {
	baseModel

	ID        string            `bun:"id,pk,type:varchar(36)"`
	Name      string            `bun:"name,type:varchar(255)"`
	UserID    string            `bun:"user_id,type:varchar(36)"`
	StartTime bun.NullTime      `bun:"start_time,nullzero"`
	EndTime   bun.NullTime      `bun:"end_time,nullzero"`
	Status    models.GameStatus `bun:"-"`

	Teams     []Team                  `bun:"rel:has-many,join:id=instance_id"`
	Locations []Location              `bun:"rel:has-many,join:id=instance_id"`
	Settings  models.InstanceSettings `bun:"rel:has-one,join:id=instance_id"`
}

// GetStatus returns the status of the instance
func (i *Instance) GetStatus() models.GameStatus {
	// If the start time is null, the game is closed
	if i.StartTime.Time.IsZero() {
		return models.Closed
	}

	// If the start time is in the future, the game is scheduled
	if i.StartTime.Time.UTC().After(time.Now().UTC()) {
		return models.Scheduled
	}

	// If the end time is in the past, the game is closed
	if !i.EndTime.Time.IsZero() && i.EndTime.Time.Before(time.Now().UTC()) {
		return models.Closed
	}

	// If the start time is in the past, the game is active
	return models.Active

}
