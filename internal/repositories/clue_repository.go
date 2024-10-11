package repositories

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type ClueRepository interface {
	// Save saves or updates a clue in the database
	Save(ctx context.Context, c *models.Clue) error
	// Delete removes the clue from the database
	Delete(ctx context.Context, c *models.Clue) error
	// FindCluesByLocation returns all clues for a given location
	FindCluesByLocation(ctx context.Context, locationID string) ([]models.Clue, error)
}

type clueRepository struct{}

// NewClueRepository creates a new ClueRepository
func NewClueRepository() ClueRepository {
	return &clueRepository{}
}

// Save saves or updates a clue in the database
func (r *clueRepository) Save(ctx context.Context, c *models.Clue) error {
	var err error
	if c.ID == "" {
		_, err = db.DB.NewInsert().Model(c).Exec(ctx)
	} else {
		_, err = db.DB.NewUpdate().Model(c).WherePK().Exec(ctx)
	}
	return err
}

// Delete removes the clue from the database
func (r *clueRepository) Delete(ctx context.Context, c *models.Clue) error {
	_, err := db.DB.NewDelete().Model(c).WherePK().Exec(ctx)
	return err
}

// FindCluesByLocation returns all clues for a given location
func (r *clueRepository) FindCluesByLocation(ctx context.Context, locationID string) ([]models.Clue, error) {
	var clues []models.Clue
	err := db.DB.NewSelect().Model(&clues).Where("location_id = ?", locationID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding clues by location: %w", err)
	}
	return clues, nil
}
