package services_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/migrate"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/repositories"
	"github.com/uptrace/bun"
)

func setupLocationService(t *testing.T) (services.LocationService, *bun.DB) {
	os.Setenv("DB_CONNECTION", "file::memory:?cache=shared")
	os.Setenv("DB_TYPE", "sqlite3")
	db := db.MustOpen()
	defer db.Close()

	// Create tables
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	migrate.CreateTables(logger, db)

	clueRepo := repositories.NewClueRepository(db)
	locationRepo := repositories.NewLocationRepository(db)
	markerRepo := repositories.NewMarkerRepository(db)
	blockStateRepo := repositories.NewBlockStateRepository(db)
	blockRepo := repositories.NewBlockRepository(db, blockStateRepo)
	locationService := services.NewLocationService(clueRepo, locationRepo, markerRepo, blockRepo)
	return locationService, db
}
