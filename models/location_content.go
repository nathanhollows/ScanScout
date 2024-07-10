package models

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type LocationContent struct {
	baseModel

	ID      string `bun:",pk" json:"id"`
	Content string `bun:"," json:"content"`
}

type LocationContents []LocationContent

// Save saves or updates an instance location content
func (i *LocationContent) Save(ctx context.Context) error {
	var err error
	if i.ID == "" {
		i.ID = uuid.New().String()
		_, err = db.NewInsert().Model(i).Exec(ctx)
	} else {
		_, err = db.NewUpdate().Model(i).WherePK().Exec(ctx)
	}
	if err != nil {
		log.Error(err)
	}
	return err
}

// Delete removes the location content from the database
func (i *LocationContent) Delete(ctx context.Context) error {
	_, err := db.NewDelete().Model(i).WherePK().Exec(ctx)
	return err
}
