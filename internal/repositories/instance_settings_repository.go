package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
)

type InstanceSettingsRepository interface {
	// Save new instance settings to the database
	Save(ctx context.Context, settings *models.InstanceSettings) error
	// Update updates an instance in the database
	Update(ctx context.Context, settings *models.InstanceSettings) error
}

type instanceSettingsRepository struct{}

func NewInstanceSettingsRepository() InstanceSettingsRepository {
	return &instanceSettingsRepository{}
}

func (r *instanceSettingsRepository) Save(ctx context.Context, settings *models.InstanceSettings) error {
	if settings.InstanceID == "" {
		return fmt.Errorf("instance ID is required")
	}
	settings.CreatedAt = time.Now()
	settings.UpdatedAt = time.Now()
	_, err := db.DB.NewInsert().Model(settings).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *instanceSettingsRepository) Update(ctx context.Context, settings *models.InstanceSettings) error {
	if settings.InstanceID == "" {
		return fmt.Errorf("instance ID is required")
	}
	settings.UpdatedAt = time.Now()
	_, err := db.DB.NewUpdate().Model(settings).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
