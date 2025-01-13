package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

type LocationRepository interface {
	// Create saves or updates a location
	Create(ctx context.Context, location *models.Location) error
	// Update updates a location in the database
	Update(ctx context.Context, location *models.Location) error

	// GetByID finds a location by ID
	GetByID(ctx context.Context, locationID string) (*models.Location, error)
	// FindByInstance finds a location by instance and code
	GetByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error)
	// Find all locations for an instance
	FindByInstance(ctx context.Context, instanceID string) ([]models.Location, error)
	// FindLocationsByMarkerID finds all locations by marker ID
	FindLocationsByMarkerID(ctx context.Context, markerID string) ([]models.Location, error)

	// Delete deletes a location from the database
	// Requires a transaction as related data will also need to be deleted
	Delete(ctx context.Context, tx *bun.Tx, locationID string) error
	// DeleteByInstance deletes all locations for an instance
	// Requires a transaction as related data will also need to be deleted
	DeleteByInstance(ctx context.Context, tx *bun.Tx, instanceID string) error

	// LoadRelations loads all relations for a location
	LoadRelations(ctx context.Context, location *models.Location) error
	// LoadClues loads all clues for a location
	LoadClues(ctx context.Context, location *models.Location) error
	// LoadMarker loads the marker for a location
	LoadMarker(ctx context.Context, location *models.Location) error
	// LoadInstance loads the instance for a location
	LoadInstance(ctx context.Context, location *models.Location) error
	// LoadBlocks loads the blocks for a location
	LoadBlocks(ctx context.Context, location *models.Location) error
}

type locationRepository struct {
	db *bun.DB
}

// NewClueRepository creates a new ClueRepository
func NewLocationRepository(db *bun.DB) LocationRepository {
	return &locationRepository{
		db: db,
	}
}

// Create saves or updates a location
func (r *locationRepository) Create(ctx context.Context, location *models.Location) error {
	var err error
	if location.ID == "" {
		location.ID = uuid.New().String()
		_, err = r.db.NewInsert().Model(location).Exec(ctx)
		return err
	}
	return r.Update(ctx, location)
}

// Update updates a location in the database
func (r *locationRepository) Update(ctx context.Context, location *models.Location) error {
	_, err := r.db.NewUpdate().Model(location).WherePK().Exec(ctx)
	return err
}

// GetByID finds a location by ID
func (r *locationRepository) GetByID(ctx context.Context, locationID string) (*models.Location, error) {
	var location models.Location
	err := r.db.NewSelect().Model(&location).Where("id = ?", locationID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding location: %w", err)
	}
	return &location, nil
}

// FindByInstance finds a location by instance and code
func (r *locationRepository) GetByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error) {
	var location models.Location
	err := r.db.
		NewSelect().
		Model(&location).
		Where("instance_id = ? AND marker_id = ?", instanceID, code).
		Relation("Marker").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding location by instance and code: %w", err)
	}
	return &location, nil
}

// Find all locations for an instance
func (r *locationRepository) FindByInstance(ctx context.Context, instanceID string) ([]models.Location, error) {
	var locations []models.Location
	err := r.db.
		NewSelect().
		Model(&locations).
		Where("instance_id = ?", instanceID).
		Relation("Marker").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all locations: %w", err)
	}
	return locations, nil
}

// FindLocationsByMarkerID finds all locations for a given marker
func (r *locationRepository) FindLocationsByMarkerID(ctx context.Context, markerID string) ([]models.Location, error) {
	var locations []models.Location
	err := r.db.NewSelect().Model(&locations).Where("marker_id = ?", markerID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding locations by marker ID: %w", err)
	}
	return locations, nil
}

// Delete deletes a location from the database
func (r *locationRepository) Delete(ctx context.Context, tx *bun.Tx, locationID string) error {
	_, err := tx.NewDelete().Model(&models.Location{ID: locationID}).WherePK().Exec(ctx)
	return err
}

// DeleteByInstance deletes all locations for an instance
func (r *locationRepository) DeleteByInstance(ctx context.Context, tx *bun.Tx, instanceID string) error {
	_, err := tx.NewDelete().Model(&models.Location{}).Where("instance_id = ?", instanceID).Exec(ctx)
	return err
}

// DeleteByInstanceWithTransaction deletes all locations for an instance with a transaction
func (r *locationRepository) DeleteByInstanceWithTransaction(ctx context.Context, tx *bun.Tx, instanceID string) error {
	_, err := tx.NewDelete().Model(&models.Location{}).Where("instance_id = ?", instanceID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting locations by instance ID: %w", err)
	}
	return nil
}

// LoadRelations loads all relations for a location
func (r *locationRepository) LoadRelations(ctx context.Context, location *models.Location) error {
	err := r.db.NewSelect().
		Model(location).
		Relation("Clues").
		Relation("Blocks").
		Relation("Instance").
		Relation("Marker").
		WherePK().
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("loading relations for location: %w", err)
	}
	return nil
}

// LoadClues loads all clues for a location
func (r *locationRepository) LoadClues(ctx context.Context, location *models.Location) error {
	err := r.db.NewSelect().
		Model(&location.Clues).
		Where("location_id = ?", location.ID).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("loading clues for location: %w", err)
	}
	return nil
}

// LoadMarker loads the marker for a location
func (r *locationRepository) LoadMarker(ctx context.Context, location *models.Location) error {
	err := r.db.NewSelect().
		Model(location).
		Relation("Marker").
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("loading marker for location: %w", err)
	}
	return nil
}

// LoadInstance loads the instance for a location
func (r *locationRepository) LoadInstance(ctx context.Context, location *models.Location) error {
	err := r.db.NewSelect().
		Model(location).
		Relation("Instance").
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("loading instance for location: %w", err)
	}
	return nil
}

// LoadBlocks loads the blocks for a location
func (r *locationRepository) LoadBlocks(ctx context.Context, location *models.Location) error {
	err := r.db.NewSelect().
		Model(location).
		Relation("Blocks").
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("loading blocks for location: %w", err)
	}
	return nil
}
