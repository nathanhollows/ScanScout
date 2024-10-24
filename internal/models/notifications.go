package models

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/db"
)

type Notification struct {
	baseModel

	ID        string `bun:"id,pk,notnull"`
	Content   string `bun:"content,type:varchar(255)"`
	Type      string `bun:"type,type:varchar(255)"`
	TeamCode  string `bun:"team_code,type:varchar(36)"`
	Dismissed bool   `bun:"dismissed,type:bool"`
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
func FindNotificationsByTeamCode(ctx context.Context, teamCode string) ([]Notification, error) {
	var notifications []Notification
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
