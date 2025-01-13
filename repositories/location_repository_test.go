package repositories_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/repositories"
)

func setupLocationRepo(t *testing.T) (repositories.LocationRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	locationRepo := repositories.NewLocationRepository(db)
	return locationRepo, cleanup
}
