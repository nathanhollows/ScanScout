package models

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Clue struct {
	baseModel

	ID         string `bun:",pk,type:varchar(36)" json:"id"`
	InstanceID string `bun:",notnull" json:"instance_id"`
	LocationID string `bun:",notnull" json:"location_id"`
	Content    string `bun:",type:text" json:"content"`
	Type       string `bun:",type:varchar(255)" json:"type"` // content, clue, map

	Instance Instance `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Location Location `bun:"rel:has-one,join:location_id=id" json:"location"`
}

type Clues []Clue

// Save saves or updates a clue
func (c *Clue) Save(ctx context.Context) error {
	var err error
	if c.ID == "" {
		c.ID = uuid.New().String()
		_, err = db.DB.NewInsert().Model(c).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(c).WherePK().Exec(ctx)
	}
	if err != nil {
		log.Error(err)
	}
	return err
}

// Delete removes the clue from the database
func (c *Clue) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(c).WherePK().Exec(ctx)
	return err
}

// FindCluesByLocation returns all clues for a given location
func FindCluesByLocation(ctx context.Context, locationID string) ([]*Clue, error) {
	var clues []*Clue
	err := db.DB.NewSelect().Model(&clues).Where("location_id = ?", locationID).Scan(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return clues, nil
}
