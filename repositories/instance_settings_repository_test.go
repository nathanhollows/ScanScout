package repositories_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/repositories"
)

func setupInstanceSettingsRepo(t *testing.T) (repositories.InstanceSettingsRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	instanceSettingsRepo := repositories.NewInstanceSettingsRepository(db)
	return instanceSettingsRepo, cleanup
}
