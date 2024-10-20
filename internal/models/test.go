package models

import (
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
