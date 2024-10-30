package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	internalModels "github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

type InstanceRepository interface {
	// Save saves an instance to the database
	Save(ctx context.Context, instance *internalModels.Instance) error
	// Update updates an instance in the database
	Update(ctx context.Context, instance *internalModels.Instance) error
	// Delete deletes an instance from the database
	Delete(ctx context.Context, instanceID string) error
}

type instanceRepository struct{}

func NewInstanceRepository() InstanceRepository {
	return &instanceRepository{}
}

func (r *instanceRepository) Save(ctx context.Context, instance *internalModels.Instance) error {
	if instance.ID == "" {
		instance.ID = uuid.New().String()
	}
	_, err := db.DB.NewInsert().Model(instance).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *instanceRepository) Update(ctx context.Context, instance *internalModels.Instance) error {
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

	// Delete instance
	_, err = tx.NewDelete().Model(&internalModels.Instance{}).Where("id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete teams
	_, err = tx.NewDelete().Model(&internalModels.Team{}).Where("instance_id = ?", instanceID).Exec(ctx)
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
	_, err = tx.NewDelete().Model(&internalModels.CheckIn{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete locations
	_, err = tx.NewDelete().Model(&internalModels.Location{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
