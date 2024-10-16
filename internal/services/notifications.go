package services

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
)

type NotificationService interface {
	SendNotification(ctx context.Context, teamCode string, content string) (models.Notification, error)
	SendNotificationToAll(ctx context.Context, team []models.Team, content string) error
	GetNotifications(ctx context.Context, teamCode string) (models.Notifications, error)
	DismissNotification(ctx context.Context, notificationID string) error
}

type notificationService struct{}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

// SendNotification sends a notification to a team
func (s *notificationService) SendNotification(ctx context.Context, teamCode string, content string) (models.Notification, error) {
	notification := models.Notification{
		TeamCode: teamCode,
		Content:  content,
	}
	err := notification.Save(ctx)
	return notification, err
}

// SendNotificationToAll sends a notification to all teams
func (s *notificationService) SendNotificationToAll(ctx context.Context, team []models.Team, content string) error {
	if len(team) == 0 {
		return fmt.Errorf("no teams to send notification to")
	}
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	for _, t := range team {
		if t.HasStarted {
			_, err := s.SendNotification(ctx, t.Code, content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetNotifications retrieves all notifications for a team
func (s *notificationService) GetNotifications(ctx context.Context, teamCode string) (models.Notifications, error) {
	return models.FindNotificationsByTeamCode(ctx, teamCode)
}

// DismissNotification marks a notification as dismissed
func (s *notificationService) DismissNotification(ctx context.Context, notificationID string) error {
	notification, err := models.FindNotificationByID(ctx, notificationID)
	if err != nil {
		return err
	}
	return notification.Dismiss(ctx)
}
