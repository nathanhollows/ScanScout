package models

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/nathanhollows/Rapua/db"
)

func SetupTestDB(t *testing.T) (cleanup func()) {
	t.Helper() // Mark function as helper for better test output

	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db.MustOpen()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Create tables
	CreateTables(logger)

	return func() {
		db.DB.Close()
	}
}

func CreateTables(logger *slog.Logger) {
	var models = []interface{}{
		(*Notification)(nil),
		(*InstanceSettings)(nil),
		(*Block)(nil),
		(*TeamBlockState)(nil),
		(*Location)(nil),
		(*Clue)(nil),
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
