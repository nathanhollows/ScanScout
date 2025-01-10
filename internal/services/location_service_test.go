package services_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/repositories"
)

func setupLocationService(t *testing.T) (services.LocationService, func()) {
	t.Helper()
	dbc, cleanup := setupDB(t)

	transactor := db.NewTransactor(dbc)

	clueRepo := repositories.NewClueRepository(dbc)
	locationRepo := repositories.NewLocationRepository(dbc)
	markerRepo := repositories.NewMarkerRepository(dbc)
	blockStateRepo := repositories.NewBlockStateRepository(dbc)
	blockRepo := repositories.NewBlockRepository(dbc, blockStateRepo)
	locationService := services.NewLocationService(transactor, clueRepo, locationRepo, markerRepo, blockRepo)
	return locationService, cleanup
}
