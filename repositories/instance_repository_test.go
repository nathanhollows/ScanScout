package repositories_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/repositories"
)

func setupInstanceRepo(t *testing.T) (repositories.InstanceRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	instanceRepo := repositories.NewInstanceRepository(db)
	return instanceRepo, cleanup
}
