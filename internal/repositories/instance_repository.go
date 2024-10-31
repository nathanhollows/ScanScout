package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

type InstanceRepository interface {
	// Save saves an instance to the database
	Save(ctx context.Context, instance *models.Instance) error
	// Update updates an instance in the database
	Update(ctx context.Context, instance *models.Instance) error
	// Delete deletes an instance from the database
	Delete(ctx context.Context, instanceID string) error
	// FindByID finds an instance by ID
	FindByID(ctx context.Context, id string) (*models.Instance, error)
}

type instanceRepository struct{}

func NewInstanceRepository() InstanceRepository {
	return &instanceRepository{}
}

func (r *instanceRepository) Save(ctx context.Context, instance *models.Instance) error {
	if instance.ID == "" {
		instance.ID = uuid.New().String()
	}
	_, err := db.DB.NewInsert().Model(instance).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *instanceRepository) Update(ctx context.Context, instance *models.Instance) error {
	if instance.ID == "" {
		return fmt.Errorf("ID is required")
	}
	_, err := db.DB.NewUpdate().Model(instance).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Ensure all related data is deleted
func (r *instanceRepository) Delete(ctx context.Context, instanceID string) error {
	type Tx struct {
		*sql.Tx
		db *bun.DB
	}

	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Delete settings
	_, err = tx.NewDelete().Model(&models.InstanceSettings{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete instance
	_, err = tx.NewDelete().Model(&models.Instance{}).Where("id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete teams
	_, err = tx.NewDelete().Model(&models.Team{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete Clues
	_, err = tx.NewDelete().Model(&models.Clue{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete CheckIns
	_, err = tx.NewDelete().Model(&models.CheckIn{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete locations
	_, err = tx.NewDelete().Model(&models.Location{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *instanceRepository) FindByID(ctx context.Context, id string) (*models.Instance, error) {
	instance := &models.Instance{}
	err := db.DB.NewSelect().
		Model(instance).
		Where("id = ?", id).
		Relation("Locations").
		Relation("Settings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
