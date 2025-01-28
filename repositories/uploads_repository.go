package repositories

import (
	"context"

	"github.com/uptrace/bun"
)

type UploadRepository interface {
	// Create a new upload in the database
	Create(ctx context.Context) error

	// Delete
	Delete(ctx context.Context, id string) error
}

type uploadRepository struct {
	db *bun.DB
}

func NewUploadRepository(db *bun.DB) UploadRepository {
	return &uploadRepository{
		db: db,
	}
}

// Create saves or updates a upload in the database
func (r *uploadRepository) Create(ctx context.Context) error {
	return nil
}

// Delete deletes a upload from the database
func (r *uploadRepository) Delete(ctx context.Context, id string) error {
	return nil
}
