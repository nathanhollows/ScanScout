package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/models"
)

type NotificationRepository interface {
	Save(context.Context, *models.Notification) error
}

type notificationRepository struct{}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

func (r *notificationRepository) Save(ctx context.Context, notification *models.Notification) error {
	// Validate the notification
	if notification.Content == "" {
		return fmt.Errorf("message is required")
	}
	if notification.TeamCode == "" {
		return fmt.Errorf("team_code is required")
	}

	// Generate a new ID if one doesn't exist
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}

	// Save the notification
	_, err := db.DB.NewInsert().Model(notification).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
