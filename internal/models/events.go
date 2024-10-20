package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Event struct {
	baseModel

	ID         string `bun:"id,pk,type:varchar(36)"`
	InstanceID string `bun:"instance_id,notnull"`
	Type       string `bun:"type,type:varchar(255)"`
	LocationID string `bun:"location_id,notnull"`
	Points     int    `bun:"points,notnull"`
	Active     bool   `bun:"active,notnull"`

	Instance Instance `bun:"rel:has-one,join:instance_id=id"`
	Location Location `bun:"rel:has-one,join:location_id=id"`
}

// Save saves or updates an event
func (e *Event) Save(ctx context.Context) error {
	var err error
	if e.ID == "" {
		e.ID = uuid.New().String()
		_, err = db.DB.NewInsert().Model(e).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(e).WherePK().Exec(ctx)
	}
	return err
}

// Delete removes the event from the database
func (e *Event) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(e).WherePK().Exec(ctx)
	return err
}

// FindEventByID finds an event by its ID
func FindEventByID(ctx context.Context, eventID string) (*Event, error) {
	event := &Event{}
	err := db.DB.NewSelect().Model(event).Where("id = ?", eventID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return event, nil
}
