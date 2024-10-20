package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type LocationRepository interface {
	// FindLocation finds a location by ID
	FindLocation(ctx context.Context, locationID string) (*models.Location, error)
	FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error)
	// Update updates a location in the database
	Update(ctx context.Context, location *models.Location) error
	// Save saves or updates a location
	Save(ctx context.Context, location *models.Location) error
	// LogCheckIn logs a check in for a team at a location
	LogCheckIn(ctx context.Context, team *models.Team, location *models.Location, mustCheckOut bool, validationRequired bool) (models.Scan, error)
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

// LogCheckIn logs a check in for a team at a location
func (r *locationRepository) LogCheckIn(ctx context.Context, team *models.Team, location *models.Location, mustCheckOut bool, validationRequired bool) (models.Scan, error) {
	scan := &models.Scan{
		TeamID:          team.Code,
		LocationID:      location.ID,
		TimeIn:          time.Now().UTC(),
		MustScanOut:     mustCheckOut,
		Points:          location.Points,
		BlocksCompleted: !validationRequired,
	}
	err := scan.Save(ctx)
	if err != nil {
		return models.Scan{}, fmt.Errorf("saving scan: %w", err)
	}

	location.CurrentCount++
	location.TotalVisits++
	err = r.Save(ctx, location)
	if err != nil {
		return models.Scan{}, fmt.Errorf("saving location: %w", err)
	}

	return *scan, nil
}
