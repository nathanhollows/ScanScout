package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Event struct {
	baseModel

	ID         string `bun:",pk,type:varchar(36)" json:"id"`
	InstanceID string `bun:",notnull" json:"instance_id"`
	Type       string `bun:",type:varchar(255)" json:"type"`
	LocationID string `bun:",notnull" json:"location_id"`
	Points     int    `bun:",notnull" json:"points"`
	Active     bool   `bun:",notnull" json:"active"`

	Instance Instance `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Location Location `bun:"rel:has-one,join:location_id=id" json:"location"`
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
