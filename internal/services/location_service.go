package services

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
)

type LocationService interface {
	FindLocation(ctx context.Context, locationID string) (*models.Location, error)
	FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error)
	LoadCluesForLocation(ctx context.Context, location *models.Location) error
	LoadCluesForLocations(ctx context.Context, locations *[]models.Location) error
	LoadRelations(ctx context.Context, location *models.Location) error
	IncrementVisitorStats(ctx context.Context, location *models.Location) error
	UpdateCoords(ctx context.Context, location *models.Location, lat, lng float64) error
	UpdateName(ctx context.Context, location *models.Location, name string) error
	UpdateLocation(ctx context.Context, location *models.Location, data LocationUpdateData) error
	CreateLocation(ctx context.Context, instanceID, name string, lat, lng float64, points int) (models.Location, error)
	CreateMarker(ctx context.Context, name string, lat, lng float64) (models.Marker, error)
	DuplicateLocation(ctx context.Context, location *models.Location, newInstanceID string) (models.Location, error)
	DeleteLocation(ctx context.Context, locationID string) error
}

type locationService struct {
	locationRepo repositories.LocationRepository
	clueRepo     repositories.ClueRepository
	markerRepo   repositories.MarkerRepository
	blockRepo    repositories.BlockRepository
}

// NewLocationService creates a new instance of LocationService
func NewLocationService(clueRepo repositories.ClueRepository) LocationService {
	return locationService{
		clueRepo:     clueRepo,
		locationRepo: repositories.NewLocationRepository(),
		markerRepo:   repositories.NewMarkerRepository(),
		blockRepo:    repositories.NewBlockRepository(),
	}
}

// FindLocation finds a location by ID
func (s locationService) FindLocation(ctx context.Context, locationID string) (*models.Location, error) {
	location, err := s.locationRepo.FindLocation(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("finding location: %v", err)
	}
	return location, nil
}

// FindLocationByInstanceAndCode finds a location by instance and code
func (s locationService) FindLocationByInstanceAndCode(ctx context.Context, instanceID string, code string) (*models.Location, error) {
	location, err := s.locationRepo.FindLocationByInstanceAndCode(ctx, instanceID, code)
	if err != nil {
		return nil, fmt.Errorf("finding location by instance and code: %v", err)
	}
	return location, nil
}

// LoadCluesForLocation loads the clues for a specific location if they are not already loaded
func (s locationService) LoadCluesForLocation(ctx context.Context, location *models.Location) error {
	if len(location.Clues) == 0 {
		clues, err := s.clueRepo.FindCluesByLocation(ctx, location.ID)
		if err != nil {
			return fmt.Errorf("finding clues: %v", err)
		}
		location.Clues = clues
	}
	return nil
}

// LoadCluesForLocations loads the clues for all given locations if they are not already loaded
func (s locationService) LoadCluesForLocations(ctx context.Context, locations *[]models.Location) error {
	for i := range *locations {
		err := s.LoadCluesForLocation(ctx, &(*locations)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadRelations loads the related data for a location
func (s locationService) LoadRelations(ctx context.Context, location *models.Location) error {
	err := s.locationRepo.LoadRelations(ctx, location)
	if err != nil {
		return fmt.Errorf("loading relations: %v", err)
	}
	return nil
}

// Update visitor stats for a location
func (s locationService) IncrementVisitorStats(ctx context.Context, location *models.Location) error {
	location.CurrentCount++
	location.TotalVisits++
	return s.locationRepo.Update(ctx, location)
}

// UpdateCoords updates the coordinates for a location
func (s locationService) UpdateCoords(ctx context.Context, location *models.Location, lat, lng float64) error {
	location.Marker.Lat = lat
	location.Marker.Lng = lng
	return s.markerRepo.Update(ctx, &location.Marker)
}

// UpdateName updates the name of a location
func (s locationService) UpdateName(ctx context.Context, location *models.Location, name string) error {
	location.Name = name
	return s.locationRepo.Update(ctx, location)
}

func (s locationService) UpdateLocation(ctx context.Context, location *models.Location, data LocationUpdateData) error {
	if location.Marker.Code == "" {
		s.locationRepo.LoadMarker(ctx, location)
	}

	// Set up the marker data
	update := false

	if data.Name != "" && data.Name != location.Marker.Name {
		location.Marker.Name = data.Name
		update = true
	}

	if data.Latitude >= -90 && data.Latitude <= 90 && data.Latitude != location.Marker.Lat {
		location.Marker.Lat = data.Latitude
		update = true
	}

	if data.Longitude >= -180 && data.Longitude <= 180 && data.Longitude != location.Marker.Lng {
		location.Marker.Lng = data.Longitude
		update = true
	}

	// To avoid updating markers that other games are using, we need to check if the marker is shared
	shared, err := s.markerRepo.IsShared(ctx, location.Marker.Code)
	if err != nil {
		return fmt.Errorf("checking if marker is shared: %v", err)
	}

	if shared && update {
		newMarker, err := s.CreateMarker(ctx, location.Marker.Name, location.Marker.Lat, location.Marker.Lng)
		if err != nil {
			return fmt.Errorf("creating new marker: %v", err)
		}
		location.MarkerID = newMarker.Code
	} else if update {
		err := s.markerRepo.Update(ctx, &location.Marker)
		if err != nil {
			return fmt.Errorf("updating marker: %v", err)
		}
	}

	// Now that the marker is updated, we can update the location
	// We'll assume if the marker was updated a new one was created

	if data.Points >= 0 && data.Points != location.Points {
		location.Points = data.Points
		update = true
	}

	if data.Name != "" && data.Name != location.Name {
		location.Name = data.Name
		update = true
	}

	if update {
		err := s.locationRepo.Update(ctx, location)
		if err != nil {
			return fmt.Errorf("updating location: %v", err)
		}
	}

	return nil

}

// CreateLocation creates a new location
func (s locationService) CreateLocation(ctx context.Context, instanceID, name string, lat, lng float64, points int) (models.Location, error) {
	// Create the marker
	marker, err := s.CreateMarker(ctx, name, lat, lng)
	if err != nil {
		return models.Location{}, fmt.Errorf("creating marker: %v", err)
	}

	location := models.Location{
		Name:       name,
		InstanceID: instanceID,
		MarkerID:   marker.Code,
		Points:     points,
	}
	err = s.locationRepo.Save(ctx, &location)
	if err != nil {
		return models.Location{}, fmt.Errorf("saving location: %v", err)
	}

	return location, nil
}

// CreateMarker creates a new marker
func (s locationService) CreateMarker(ctx context.Context, name string, lat, lng float64) (models.Marker, error) {
	marker := models.Marker{
		Name: name,
		Lat:  lat,
		Lng:  lng,
	}
	err := s.markerRepo.Save(ctx, &marker)
	if err != nil {
		return models.Marker{}, fmt.Errorf("saving marker: %v", err)
	}
	return marker, nil
}

// DuplicateLocation duplicates a location
func (s locationService) DuplicateLocation(ctx context.Context, location *models.Location, newInstanceID string) (models.Location, error) {
	newLocation := *location
	newLocation.ID = ""
	newLocation.InstanceID = newInstanceID
	err := s.locationRepo.Save(ctx, &newLocation)
	if err != nil {
		return models.Location{}, fmt.Errorf("saving location: %v", err)
	}
	return newLocation, nil
}

// DeleteLocation deletes a location
// It also deletes all related clues, blocks, and scans
// If the marker is not used by any other locations, it is also deleted
func (s locationService) DeleteLocation(ctx context.Context, locationID string) error {
	location, err := s.locationRepo.FindLocation(ctx, locationID)
	if err != nil {
		return fmt.Errorf("finding location: %v", err)
	}

	// Delete all related clues
	err = s.clueRepo.DeleteByLocationID(ctx, locationID)
	if err != nil {
		return fmt.Errorf("deleting clues: %v", err)
	}

	// Delete all related blocks
	err = s.blockRepo.DeleteByLocationID(ctx, locationID)
	if err != nil {
		return fmt.Errorf("deleting blocks: %v", err)
	}

	// Delete the location
	err = s.locationRepo.Delete(ctx, locationID)
	if err != nil {
		return fmt.Errorf("deleting location: %v", err)
	}

	// Delete the marker if it is not used by any other locations
	locations, err := s.locationRepo.FindLocationsByMarkerID(ctx, location.MarkerID)
	if err != nil {
		return fmt.Errorf("finding locations by marker: %v", err)
	}
	if len(locations) == 0 {
		err = s.markerRepo.Delete(ctx, location.MarkerID)
		if err != nil {
			return fmt.Errorf("deleting marker: %v", err)
		}
	}

	return nil
}
