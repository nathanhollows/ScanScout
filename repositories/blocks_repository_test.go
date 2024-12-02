package repositories_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/migrate"
	"github.com/nathanhollows/Rapua/repositories"
)

func setupBlockRepo(t *testing.T) (repositories.BlockRepository, func()) {
	t.Helper()
	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db := db.MustOpen()

	// Create tables
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	migrate.CreateTables(logger, db)

	blockStateRepo := repositories.NewBlockStateRepository(db)
	blockRepo := repositories.NewBlockRepository(db, blockStateRepo)
	return blockRepo, func() {
		db.Close()
	}
}
