package models

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
)

func CreateTables(logger *slog.Logger) {
	var models = []interface{}{
		(*models.Notification)(nil),
		(*models.InstanceSettings)(nil),
		(*models.Block)(nil),
		(*models.TeamBlockState)(nil),
		(*Location)(nil),
		(*models.Clue)(nil),
		(*Team)(nil),
		(*Marker)(nil),
		(*CheckIn)(nil),
		(*Instance)(nil),
		(*User)(nil),
	}

	for _, model := range models {
		_, err := db.DB.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}

type baseModel struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
