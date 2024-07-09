package models

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type InstanceLocationContent struct {
	baseModel

	ID      string `bun:",pk" json:"id"`
	Content string `bun:"," json:"content"`
}

type InstanceLocationContents []InstanceLocationContent

// Save saves or updates an instance location content
func (i *InstanceLocationContent) Save(ctx context.Context) error {
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
