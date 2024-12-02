package services_test

import (
	"context"
	"testing"

	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/models"
	"github.com/stretchr/testify/assert"
	"github.com/nathanhollows/Rapua/repositories"
)

func setupLocationService(t *testing.T) (services.LocationService, func()) {
	cleanup := models.SetupTestDB(t)
	clueRepo := repositories.NewClueRepository()
	locationService := services.NewLocationService(clueRepo)
	return locationService, cleanup
}

func TestLocationService_LoadCluesForLocation(t *testing.T) {
	locationService, cleanup := setupLocationService(t)
	defer cleanup()

	ctx := context.Background()
	locationID := "location-1"

	// Save some clues for testing
	clueRepo := repositories.NewClueRepository()
	clue1 := &models.Clue{
		InstanceID: "instance-1",
		LocationID: locationID,
		Content:    "Clue 1",
	}
	clue2 := &models.Clue{
		InstanceID: "instance-2",
		LocationID: locationID,
		Content:    "Clue 2",
	}
	err := clueRepo.Save(ctx, clue1)
	assert.NoError(t, err, "expected no error when saving clue 1")
	err = clueRepo.Save(ctx, clue2)
	assert.NoError(t, err, "expected no error when saving clue 2")

	// Create location and load clues
	location := &models.Location{ID: locationID}
	err = locationService.LoadCluesForLocation(ctx, location)

	assert.NoError(t, err, "expected no error when loading clues for location")
	assert.Len(t, location.Clues, 2, "expected two clues to be loaded")
}
