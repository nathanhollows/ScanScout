package models

import (
	"context"
	"testing"
)

func TestNotifcation_Save(t *testing.T) {
	SetupTestDB(t)

	tests := []struct {
		name    string
		notif   *Notification
		wantErr bool
	}{
		{
			name: "Valid Notification",
			notif: &Notification{
				Content:   "Hello, World!",
				Type:      "info",
				TeamCode:  "1",
				Dismissed: false,
			},
			wantErr: false,
		},
		{
			name:    "Invalid Notification",
			notif:   &Notification{},
			wantErr: true,
		},
		{
			name: "Invalid Notification - Missing Notification",
			notif: &Notification{
				Type:      "info",
				TeamCode:  "1",
				Dismissed: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notif.Save(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Notification.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotification_Update(t *testing.T) {
	SetupTestDB(t)

	tests := []struct {
		name    string
		notif   *Notification
		wantErr bool
	}{
		{
			name: "Valid Notification",
			notif: &Notification{
				Content:   "Hello, World!",
				Type:      "info",
				TeamCode:  "1",
				Dismissed: false,
			},
			wantErr: false,
		},
		{
			name:    "Invalid Notification",
			notif:   &Notification{},
			wantErr: true,
		},
		{
			name: "Invalid Notification - Missing ID",
			notif: &Notification{
				Content:   "Hello, World!",
				Type:      "info",
				TeamCode:  "1",
				Dismissed: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save the notification if it's valid
			// This doesn't test the Save function, but it's required for the Update function
			if tt.wantErr == false {
				tt.notif.Save(context.Background())
			}

			err := tt.notif.Update(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Notification.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotification_Delete(t *testing.T) {
	SetupTestDB(t)

	tests := []struct {
		name    string
		notif   *Notification
		wantErr bool
	}{
		{
			name: "Valid Notification",
			notif: &Notification{
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
			tt.notif.Save(context.Background())
			err := tt.notif.Delete(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Notification.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFindNotificationByID tests the FindNotificationByID function
func TestFindNotificationByID(t *testing.T) {
	SetupTestDB(t)

	tests := []struct {
		name    string
		notif   *Notification
		wantErr bool
	}{
		{
			name: "Valid Notification",
			notif: &Notification{
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
			tt.notif.Save(context.Background())
			_, err := FindNotificationByID(context.Background(), tt.notif.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindNotificationByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDismissNotification tests the DismissNotification function
func TestDismissNotification(t *testing.T) {
	SetupTestDB(t)

	tests := []struct {
		name    string
		notif   *Notification
		wantErr bool
	}{
		{
			name: "Valid Notification",
			notif: &Notification{
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
			tt.notif.Save(context.Background())
			err := tt.notif.Dismiss(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("DismissNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
