package services_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nathanhollows/Rapua/v3/db"
	"github.com/nathanhollows/Rapua/v3/internal/services"
	"github.com/nathanhollows/Rapua/v3/repositories"
	"github.com/stretchr/testify/assert"
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

func TestLocationService_DuplicateLocation(t *testing.T) {
	service, cleanup := setupLocationService(t)
	defer cleanup()

	blockService, blockCleanup := setupBlocksService(t)
	defer blockCleanup()

	t.Run("Duplicate location", func(t *testing.T) {
		// Create a location
		location, err := service.CreateLocation(
			context.Background(),
			gofakeit.UUID(),
			gofakeit.Name(),
			gofakeit.Latitude(),
			gofakeit.Longitude(),
			gofakeit.Number(0, 100))
		assert.NoError(t, err)

		// Create a block
		_, err = blockService.NewBlock(context.Background(), location.ID, "image")
		assert.NoError(t, err)

		// Duplicate the location
		newLocation, err := service.DuplicateLocation(context.Background(), location, gofakeit.UUID())
		assert.NoError(t, err)
		assert.NotEqual(t, location.ID, newLocation.ID)

		// Check that the location was duplicated
		checkLocation, err := service.GetByID(context.Background(), newLocation.ID)
		assert.NoError(t, err)
		assert.NotNil(t, checkLocation)

		// Check that the blocks were duplicated
		blocks, err := blockService.FindByLocationID(context.Background(), newLocation.ID)
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)

	})
}
