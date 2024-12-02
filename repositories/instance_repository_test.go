package repositories_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/migrate"
	"github.com/nathanhollows/Rapua/repositories"
)

func setupInstanceRepo(t *testing.T) (repositories.InstanceRepository, func()) {
	t.Helper()
	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db := db.MustOpen()

	// Create tables
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	migrate.CreateTables(logger, db)

	instanceRepo := repositories.NewInstanceRepository(db)
	return instanceRepo, func() {
		db.Close()
	}
}
