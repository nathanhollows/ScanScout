package migrate

import (
	"context"
	"log"
	"log/slog"

	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

func CreateTables(logger *slog.Logger, db *bun.DB) {
	var models = []interface{}{
		(*models.Notification)(nil),
		(*models.InstanceSettings)(nil),
		(*models.Block)(nil),
		(*models.TeamBlockState)(nil),
		(*models.Location)(nil),
		(*models.Clue)(nil),
		(*models.Team)(nil),
		(*models.Marker)(nil),
		(*models.CheckIn)(nil),
		(*models.Instance)(nil),
		(*models.User)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}
