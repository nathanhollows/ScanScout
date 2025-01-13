package repositories_test

import (
	"testing"

	"github.com/nathanhollows/Rapua/repositories"
)

func setupNotificationRepo(t *testing.T) (repositories.NotificationRepository, func()) {
	t.Helper()
	db, cleanup := setupDB(t)

	notificationRepo := repositories.NewNotificationRepository(db)
	return notificationRepo, cleanup
}
