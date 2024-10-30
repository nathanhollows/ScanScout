package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
)

type ClueRepository interface {
	// Save saves or updates a clue in the database
	Save(ctx context.Context, c *models.Clue) error
	// Delete removes the clue from the database
	Delete(ctx context.Context, clueID string) error
	// DeleteByLocationID removes all clues for a location
	DeleteByLocationID(ctx context.Context, locationID string) error
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
	if c.InstanceID == "" || c.LocationID == "" {
		return fmt.Errorf("instance ID and location ID must be set")
	}
	var err error
	if c.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generating UUID: %w", err)
		}
		c.ID = id.String()
	}
	_, err = db.DB.NewInsert().Model(c).Exec(ctx)
	return err
}

// Delete removes the clue from the database
func (r *clueRepository) Delete(ctx context.Context, clueID string) error {
	_, err := db.DB.
		NewDelete().
		Model(&models.Clue{ID: clueID}).
		ForceDelete().
		WherePK().
		Exec(ctx)
	return err
}

// DeleteByLocationID removes all clues for a location
func (r *clueRepository) DeleteByLocationID(ctx context.Context, locationID string) error {
	_, err := db.DB.
		NewDelete().
		Model(&models.Clue{}).
		Where("location_id = ?", locationID).
		ForceDelete().
		Exec(ctx)
	return err
}

// FindCluesByLocation returns all clues for a given location
func (r *clueRepository) FindCluesByLocation(ctx context.Context, locationID string) ([]models.Clue, error) {
	clues := []models.Clue{}
	err := db.DB.
		NewSelect().
		Model(&clues).
		Where("location_id = ?", locationID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding clues by location: %w", err)
	}
	return clues, nil
}
