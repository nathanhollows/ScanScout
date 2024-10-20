package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
)

// Instance represents a single planned activity belonging to a user
// Instance is used to match users, teams, locations, and scans
type Instance struct {
	baseModel

	ID        string       `bun:",pk,type:varchar(36)" json:"id"`
	Name      string       `bun:",type:varchar(255)" json:"name"`
	UserID    string       `bun:",type:varchar(36)" json:"user_id"`
	StartTime bun.NullTime `bun:",nullzero" json:"start_time"`
	EndTime   bun.NullTime `bun:",nullzero" json:"end_time"`
	Status    GameStatus   `bun:"-" json:"status"`

	Teams     []Team           `bun:"rel:has-many,join:id=instance_id" json:"teams"`
	Locations Locations        `bun:"rel:has-many,join:id=instance_id" json:"instance_locations"`
	Scans     []Scan           `bun:"rel:has-many,join:id=instance_id" json:"scans"`
	Settings  InstanceSettings `bun:"rel:has-one,join:id=instance_id" json:"settings"`
}

func (i *Instance) Save(ctx context.Context) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	_, err := db.DB.NewInsert().Model(i).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Update(ctx context.Context) error {
	_, err := db.DB.NewUpdate().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Deleting an instance will cascade delete all teams, locations, and scans
func (i *Instance) Delete(ctx context.Context) error {
	// Delete teams
	for _, team := range i.Teams {
		err := team.Delete(ctx)
		if err != nil {
			return err
		}
	}

	// Delete locations
	for _, location := range i.Locations {
		err := location.Delete(ctx)
		if err != nil {
			return err
		}
	}

	// Delete scans
	for _, scan := range i.Scans {
		err := scan.Delete(ctx)
		if err != nil {
			return err
		}
	}

	_, err := db.DB.NewDelete().Model(i).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FindAllInstances finds all instances
func FindAllInstances(ctx context.Context, userID string) ([]Instance, error) {
	instances := []Instance{}
	err := db.DB.NewSelect().Model(&instances).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// FindInstanceByID finds an instance by ID
func FindInstanceByID(ctx context.Context, id string) (*Instance, error) {
	instance := &Instance{}
	err := db.DB.NewSelect().
		Model(instance).
		Where("id = ?", id).
		Relation("Locations").
		Relation("Settings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// GetStatus returns the status of the instance
func (i *Instance) GetStatus() GameStatus {
	// If the start time is null, the game is closed
	if i.StartTime.Time.IsZero() {
		return Closed
	}

	// If the start time is in the future, the game is scheduled
	if i.StartTime.Time.UTC().After(time.Now().UTC()) {
		return Scheduled
	}

	// If the end time is in the past, the game is closed
	if !i.EndTime.Time.IsZero() && i.EndTime.Time.Before(time.Now().UTC()) {
		return Closed
	}

	// If the start time is in the past, the game is active
	return Active

}

// LoadSettings loads the settings for an instance
func (i *Instance) LoadSettings(ctx context.Context) error {
	if i.Settings.InstanceID == "" {
		i.Settings = InstanceSettings{}
		err := db.DB.NewSelect().Model(&i.Settings).Where("instance_id = ?", i.ID).Scan(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadLocations loads the locations for an instance
func (i *Instance) LoadLocations(ctx context.Context) error {
	if len(i.Locations) > 0 {
		return nil
	}

	var err error
	i.Locations, err = FindAllLocations(ctx, i.ID)
	if err != nil {
		return err
	}

	return nil
}

// LoadTeams loads the teams for an instance
func (i *Instance) LoadTeams(ctx context.Context) error {
	if len(i.Teams) > 0 {
		return nil
	}

	var err error
	i.Teams, err = FindAllTeams(ctx, i.ID)
	if err != nil {
		return err
	}

	return nil
}
