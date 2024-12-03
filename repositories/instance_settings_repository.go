package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

type InstanceSettingsRepository interface {
	// Create new instance settings to the database
	Create(ctx context.Context, settings *models.InstanceSettings) error

	// Update updates an instance in the database
	Update(ctx context.Context, settings *models.InstanceSettings) error
}

type instanceSettingsRepository struct {
	db *bun.DB
}

func NewInstanceSettingsRepository(db *bun.DB) InstanceSettingsRepository {
	return &instanceSettingsRepository{
		db: db,
	}
}

func (r *instanceSettingsRepository) Create(ctx context.Context, settings *models.InstanceSettings) error {
	if settings.InstanceID == "" {
		return fmt.Errorf("instance ID is required")
	}
	settings.CreatedAt = time.Now()
	settings.UpdatedAt = time.Now()
	_, err := r.db.NewInsert().Model(settings).Exec(ctx)
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
	_, err := r.db.NewUpdate().Model(settings).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
