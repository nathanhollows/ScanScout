package repositories_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/repositories"
)

func setupMarkerRepo(t *testing.T) (repositories.MarkerRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	markerRepo := repositories.NewMarkerRepository(db)
	return markerRepo, cleanup
}
