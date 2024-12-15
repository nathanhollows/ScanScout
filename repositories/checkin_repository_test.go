package repositories_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/repositories"
)

func setupCheckinRepo(t *testing.T) (repositories.CheckInRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	checkinRepo := repositories.NewCheckInRepository(db)
	return checkinRepo, cleanup
}
