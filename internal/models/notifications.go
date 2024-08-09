package models

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type Notification struct {
	baseModel

	ID        string `bun:",pk,notnull" json:"id"`
	Content   string `bun:",type:varchar(255)" json:"content"`
	Type      string `bun:",type:varchar(255)" json:"type"`
	TeamCode  string `bun:",type:varchar(36)" json:"team_code"`
	Dismissed bool   `bun:",type:bool" json:"dismissed"`
}

type Notifications []*Notification

// Save saves a notification
func (n *Notification) Save(ctx context.Context) error {
	// Validate the notification
	if n.Content == "" {
		return fmt.Errorf("message is required")
	}
	if n.TeamCode == "" {
		return fmt.Errorf("team_id is required")
	}

	// Generate a new ID if one doesn't exist
	if n.ID == "" {
		n.ID = uuid.New().String()
	}

	// Save the notification
	_, err := db.DB.NewInsert().Model(n).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a notification
func (n *Notification) Update(ctx context.Context) error {
	if n.ID == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := db.DB.NewUpdate().Model(n).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a notification
func (n *Notification) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(n).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FindNotificationByID finds a notification by ID
func FindNotificationByID(ctx context.Context, id string) (*Notification, error) {
	notication := &Notification{}
	err := db.DB.NewSelect().Model(notication).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return notication, nil
}

// FindNotificationsByTeamCode finds notifications by team code
func FindNotificationsByTeamCode(ctx context.Context, teamCode string) (Notifications, error) {
	var notifications Notifications
	err := db.DB.NewSelect().Model(&notifications).Where("team_code = ? AND NOT dismissed", teamCode).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("FindNotificationsByTeamCode: %w", err)
	}
	return notifications, nil
}

// Dismiss marks a notification as dismissed
func (n *Notification) Dismiss(ctx context.Context) error {
	n.Dismissed = true
	return n.Update(ctx)
}
