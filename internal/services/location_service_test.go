package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/stretchr/testify/assert"
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
		ID:         uuid.New().String(),
		InstanceID: "instance-1",
		LocationID: locationID,
		Content:    "Clue 1",
	}
	clue2 := &models.Clue{
		ID:         uuid.New().String(),
		InstanceID: "instance-2",
		LocationID: locationID,
		Content:    "Clue 2",
	}
	clueRepo.Save(ctx, clue1)
	clueRepo.Save(ctx, clue2)

	// Create location and load clues
	location := &models.Location{ID: locationID}
	err := locationService.LoadCluesForLocation(ctx, location)

	assert.NoError(t, err, "expected no error when loading clues for location")
	assert.Len(t, location.Clues, 2, "expected two clues to be loaded")
}
