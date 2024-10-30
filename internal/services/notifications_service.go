package services

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type NotificationService interface {
	SendNotification(ctx context.Context, teamCode string, content string) (models.Notification, error)
	SendNotificationToAllTeams(ctx context.Context, instanceID string, content string) error
	GetNotifications(ctx context.Context, teamCode string) ([]models.Notification, error)
	DismissNotification(ctx context.Context, notificationID string) error
}

type notificationService struct {
	teamRepository         repositories.TeamRepository
	notificationRepository repositories.NotificationRepository
}

func NewNotificationService() NotificationService {
	return &notificationService{
		teamRepository:         repositories.NewTeamRepository(),
		notificationRepository: repositories.NewNotificationRepository(),
	}
}

// SendNotification sends a notification to a team
func (s *notificationService) SendNotification(ctx context.Context, teamCode string, content string) (models.Notification, error) {
	notification := models.Notification{
		TeamCode: teamCode,
		Content:  content,
	}

	err := s.notificationRepository.Save(ctx, &notification)
	return notification, err
}

// SendNotificationToAllTeams sends a notification to all teams
func (s *notificationService) SendNotificationToAllTeams(ctx context.Context, instanceID string, content string) error {
	teams, err := s.teamRepository.FindAll(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("error finding teams: %w", err)
	}

	if len(teams) == 0 {
		return fmt.Errorf("no teams to send notification to")
	}

	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	for _, team := range teams {
		if team.HasStarted {
			_, err := s.SendNotification(ctx, team.Code, content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetNotifications retrieves all notifications for a team
func (s *notificationService) GetNotifications(ctx context.Context, teamCode string) ([]models.Notification, error) {
	return s.notificationRepository.FindByTeamCode(ctx, teamCode)
}

// DismissNotification marks a notification as dismissed
func (s *notificationService) DismissNotification(ctx context.Context, notificationID string) error {
	notification, err := s.notificationRepository.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}
	notification.Dismissed = true
	return s.notificationRepository.Update(ctx, &notification)
}
