package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/models"
)

type LocationRepository interface {
	// FindLocation finds a location by ID
	FindLocation(ctx context.Context, locationID string) (*models.Location, error)
	FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error)
	// Update updates a location in the database
	Update(ctx context.Context, location *models.Location) error
	// Save saves or updates a location
	Save(ctx context.Context, location *models.Location) error
}

type locationRepository struct{}

// NewClueRepository creates a new ClueRepository
func NewLocationRepository() LocationRepository {
	return &locationRepository{}
}

// FindLocation finds a location by ID
func (r *locationRepository) FindLocation(ctx context.Context, locationID string) (*models.Location, error) {
	var location models.Location
	err := db.DB.NewSelect().Model(&location).Where("id = ?", locationID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding location: %w", err)
	}
	return &location, nil
}

// FindLocationByInstanceAndCode finds a location by instance and code
func (r *locationRepository) FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error) {
	var location models.Location
	err := db.DB.NewSelect().Model(&location).Where("instance_id = ? AND marker_id = ?", instanceID, code).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding location by instance and code: %w", err)
	}
	return &location, nil
}

// Update updates a location in the database
func (r *locationRepository) Update(ctx context.Context, location *models.Location) error {
	_, err := db.DB.NewUpdate().Model(location).WherePK().Exec(ctx)
	return err
}

// Save saves or updates a location
func (r *locationRepository) Save(ctx context.Context, location *models.Location) error {
	var err error
	if location.ID == "" {
		location.ID = uuid.New().String()
		_, err = db.DB.NewInsert().Model(location).Exec(ctx)
		return err
	}
	return r.Update(ctx, location)
}
