package models

import (
	"os"
	"testing"

	"github.com/nathanhollows/Rapua/pkg/db"
)

func SetupTestDB(t *testing.T) (cleanup func()) {
	t.Helper() // Mark function as helper for better test output

	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db.Connect()

	// Create tables
	CreateTables()

	return func() {
		db.DB.Close()
	}
}
