package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/models"
)

type NotificationRepository interface {
	//	Save saves a notification to the database
	Save(context.Context, *models.Notification) error
	//	Update updates a notification in the database
	Update(context.Context, *models.Notification) error
	//	Delete deletes a notification from the database
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (models.Notification, error)
	FindByTeamCode(ctx context.Context, teamCode string) ([]models.Notification, error)
}

type notificationRepository struct{}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

// Save inserts a new notification into the database
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

// Update updates an existing notification in the database
func (r *notificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	if notification.ID == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := db.DB.NewUpdate().Model(notification).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a notification from the database
func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	_, err := db.DB.NewDelete().Model(&models.Notification{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FindByID finds a notification by its ID
func (r *notificationRepository) FindByID(ctx context.Context, id string) (models.Notification, error) {
	notification := models.Notification{}
	err := db.DB.NewSelect().Model(&notification).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return notification, err
	}
	return notification, nil
}

// FindByTeamCode finds all notifications for a specific team code
func (r *notificationRepository) FindByTeamCode(ctx context.Context, teamCode string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := db.DB.NewSelect().Model(&notifications).Where("team_code = ? AND NOT dismissed", teamCode).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("FindByTeamCode: %w", err)
	}
	return notifications, nil
}
