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
	// FindLocationByInstanceAndCode finds a location by instance and code
	FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error)
	// Find all locations for an instance
	FindAllLocations(ctx context.Context, instanceID string) ([]models.Location, error)
	// Update updates a location in the database
	// FindLocationsByMarkerID(ctx context.Context, markerID string) ([]models.Location, error)
	FindLocationsByMarkerID(ctx context.Context, markerID string) ([]models.Location, error)
	Update(ctx context.Context, location *models.Location) error
	// Save saves or updates a location
	Save(ctx context.Context, location *models.Location) error
	// Delete deletes a location from the database
	Delete(ctx context.Context, locationID string) error
}

type locationRepository struct {
}

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

// Find all locations for an instance
func (r *locationRepository) FindAllLocations(ctx context.Context, instanceID string) ([]models.Location, error) {
	var locations []models.Location
	err := db.DB.NewSelect().Model(&locations).Where("instance_id = ?", instanceID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all locations: %w", err)
	}
	return locations, nil
}

// FindLocationsByMarkerID finds all locations for a given marker
func (r *locationRepository) FindLocationsByMarkerID(ctx context.Context, markerID string) ([]models.Location, error) {
	var locations []models.Location
	err := db.DB.NewSelect().Model(&locations).Where("marker_id = ?", markerID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding locations by marker ID: %w", err)
	}
	return locations, nil
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

// Delete deletes a location from the database
func (r *locationRepository) Delete(ctx context.Context, locationID string) error {
	_, err := db.DB.NewDelete().Model(&models.Location{ID: locationID}).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting location: %w", err)
	}
	return nil
}
