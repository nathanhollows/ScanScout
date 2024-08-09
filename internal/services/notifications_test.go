package services_test

import (
	"context"
	"testing"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationService_SendNotification(t *testing.T) {
	dbCleanupFunc := models.SetupTestDB(t)
	defer dbCleanupFunc()

	// Create NotificationService
	notificationService := services.NewNotificationService()

	tests := []struct {
		name     string
		teamCode string
		content  string
		wantErr  bool
	}{
		{
			name:     "Valid Notification",
			teamCode: "team1",
			content:  "This is a test notification",
			wantErr:  false,
		},
		{
			name:     "Empty TeamCode",
			teamCode: "",
			content:  "This is a test notification",
			wantErr:  true,
		},
		{
			name:     "Empty Content",
			teamCode: "team1",
			content:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send Notification
			notification, err := notificationService.SendNotification(context.Background(), tt.teamCode, tt.content)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error: %v", err)
				assert.NotZero(t, notification.ID, "notification ID should be non-zero")
				assert.Equal(t, notification.TeamCode, tt.teamCode, "team codes should match")
				assert.Equal(t, notification.Content, tt.content, "contents should match")
			}
		})
	}
}

func TestNotificationService_SendNotificationToAll(t *testing.T) {
	dbCleanupFunc := models.SetupTestDB(t)
	defer dbCleanupFunc()

	// Create NotificationService
	notificationService := services.NewNotificationService()

	tests := []struct {
		name    string
		teams   models.Teams
		content string
		wantErr bool
	}{
		{
			name:    "Valid Notification",
			teams:   models.Teams{{Code: "team1"}, {Code: "team2"}},
			content: "This is a test notification",
			wantErr: false,
		},
		{
			name:    "No Teams",
			teams:   models.Teams{},
			content: "This is a test notification",
			wantErr: true,
		},
		{
			name:    "Empty Content",
			teams:   models.Teams{{Code: "team1"}, {Code: "team2"}},
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send Notification
			err := notificationService.SendNotificationToAll(context.Background(), tt.teams, tt.content)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error: %v", err)
			}
		})
	}
}

func TestNotificationService_GetNotifications(t *testing.T) {
	dbCleanupFunc := models.SetupTestDB(t)
	defer dbCleanupFunc()

	// Create NotificationService
	notificationService := services.NewNotificationService()

	tests := []struct {
		name               string
		teamCode           string
		setupNotifications []models.Notification
		wantLen            int
		wantErr            bool
	}{
		{
			name:     "No Notifications",
			teamCode: "team1",
			wantLen:  0,
			wantErr:  false,
		},
		{
			name:     "One Notification",
			teamCode: "team1",
			setupNotifications: []models.Notification{
				{TeamCode: "team1", Content: "This is a test notification"},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:     "Multiple Notifications",
			teamCode: "team2",
			setupNotifications: []models.Notification{
				{TeamCode: "team2", Content: "First notification for team2"},
				{TeamCode: "team2", Content: "Second notification for team2"},
				{TeamCode: "team2", Content: "Third notification for team2"},
			},
			wantLen: 3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial notifications
			for _, notif := range tt.setupNotifications {
				_, err := notificationService.SendNotification(context.Background(), notif.TeamCode, notif.Content)
				require.NoError(t, err, "unexpected error: %v", err)
			}

			// Get Notifications
			notifications, err := notificationService.GetNotifications(context.Background(), tt.teamCode)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error: %v", err)
				assert.Len(t, notifications, tt.wantLen, "number of notifications should match")
			}
		})
	}
}

// TestDismissNotification tests the DismissNotification function
func TestNotificationService_DismissNotification(t *testing.T) {
	dbCleanupFunc := models.SetupTestDB(t)
	defer dbCleanupFunc()

	// Create NotificationService
	notificationService := services.NewNotificationService()

	tests := []struct {
		name    string
		notif   *models.Notification
		wantErr bool
	}{
		{
			name: "Valid Notification",
			notif: &models.Notification{
				Content:   "Hello, World!",
				Type:      "info",
				TeamCode:  "1",
				Dismissed: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save Notification
			err := tt.notif.Save(context.Background())
			require.NoError(t, err, "unexpected error: %v", err)

			// Dismiss Notification
			err = notificationService.DismissNotification(context.Background(), tt.notif.ID)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error: %v", err)

				// Check if notification is dismissed
				notif, err := models.FindNotificationByID(context.Background(), tt.notif.ID)
				require.NoError(t, err, "unexpected error: %v", err)
				assert.True(t, notif.Dismissed, "notification should be dismissed")
			}
		})
	}
}
