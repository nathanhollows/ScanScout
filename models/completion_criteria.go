package models

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type CompletionCriteria struct {
	baseModel

	ID   string `bun:",pk,type:varchar(36)" json:"id"`
	Name string `bun:",type:varchar(255)" json:"name"`
}

type CompletionCriterias []CompletionCriteria

// Save saves or updates a completion criteria
func (c *CompletionCriteria) Save(ctx context.Context) error {
	var err error
	if c.ID == "" {
		c.ID = uuid.New().String()
		_, err = db.NewInsert().Model(c).Exec(ctx)
	} else {
		_, err = db.NewUpdate().Model(c).WherePK().Exec(ctx)
	}
	if err != nil {
		log.Error(err)
	}
	return err
}
