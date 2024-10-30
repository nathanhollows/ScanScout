package repositories

import (
	"context"
	"strings"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
)

type MarkerRepository interface {
	// Create a new marker in the database
	Save(ctx context.Context, marker *models.Marker) error
	// Update a marker in the database
	Update(ctx context.Context, marker *models.Marker) error
	// Delete
	Delete(ctx context.Context, code string) error
	FindByCode(ctx context.Context, code string) (*models.Marker, error)
	UpdateCoords(ctx context.Context, marker *models.Marker, lat, lng float64) error
	// Is Shared checks if a marker is used by more than one location
	IsShared(ctx context.Context, code string) (bool, error)
}

type markerRepository struct{}

func NewMarkerRepository() MarkerRepository {
	return &markerRepository{}
}

// Save saves or updates a marker in the database
func (r *markerRepository) Save(ctx context.Context, marker *models.Marker) error {
	if marker.Code == "" {
		// TODO: Remove magic number
		marker.Code = helpers.NewCode(5)
		_, err := db.DB.NewInsert().Model(marker).Exec(ctx)
		return err
	}
	_, err := db.DB.NewUpdate().Model(marker).WherePK("code").Exec(ctx)
	return err
}

// Update updates a marker in the database
func (r *markerRepository) Update(ctx context.Context, marker *models.Marker) error {
	_, err := db.DB.
		NewUpdate().
		Model(marker).
		Column("name", "lat", "lng", "total_visits", "current_count", "avg_duration").
		WherePK("code").
		Exec(ctx)

	return err
}

// Delete deletes a marker from the database
func (r *markerRepository) Delete(ctx context.Context, markerCode string) error {
	_, err := db.DB.NewDelete().Model(&models.Marker{Code: markerCode}).WherePK().Exec(ctx)
	return err
}

// FindByCode finds a marker by its code
func (r *markerRepository) FindByCode(ctx context.Context, code string) (*models.Marker, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	var marker models.Marker
	err := db.DB.NewSelect().Model(&marker).Where("code = ?", code).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &marker, nil
}

// UpdateCoords updates the latitude and longitude of a marker in the database
func (r *markerRepository) UpdateCoords(ctx context.Context, marker *models.Marker, lat, lng float64) error {
	marker.Lat = lat
	marker.Lng = lng
	_, err := db.DB.NewUpdate().Model(marker).WherePK().Column("lat", "lng").Exec(ctx)
	return err
}

// IsShared checks if a marker is shared
func (r *markerRepository) IsShared(ctx context.Context, code string) (bool, error) {
	var count int
	count, err := db.DB.NewSelect().Model(&models.Location{}).Where("marker_id = ?", code).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 1, nil
}
