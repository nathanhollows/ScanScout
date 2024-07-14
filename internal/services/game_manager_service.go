package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type GameManagerService struct{}

func NewGameManagerService() *GameManagerService {
	return &GameManagerService{}
}

func (s *GameManagerService) CreateInstance(ctx context.Context, name string, user *models.User) (*models.Instance, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	instance := &models.Instance{
		Name:   name,
		UserID: user.ID,
	}

	if err := instance.Save(ctx); err != nil {
		return nil, err
	}

	user.CurrentInstanceID = instance.ID
	if err := user.Update(ctx); err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *GameManagerService) SwitchInstance(ctx context.Context, user *models.User, instanceID string) (*models.Instance, error) {
	instance, err := models.FindInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("instance not found")
	}

	user.CurrentInstanceID = instance.ID
	if err := user.Update(ctx); err != nil {
		return nil, err
	}

	return instance, nil
}

func (s *GameManagerService) DeleteInstance(ctx context.Context, user *models.User, instanceID, confirmName string) error {
	instance, err := models.FindInstanceByID(ctx, instanceID)
	if err != nil {
		return errors.New("instance not found")
	}

	if user.ID != instance.UserID {
		return errors.New("you do not have permission to delete this instance")
	}

	if confirmName != instance.Name {
		return errors.New("instance name does not match confirmation")
	}

	if user.CurrentInstanceID == instance.ID {
		user.CurrentInstanceID = ""
		if err := user.Update(ctx); err != nil {
			return err
		}
	}

	if err := instance.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (s *GameManagerService) AddTeams(ctx context.Context, user *models.User, count int) error {
	if count < 1 {
		return errors.New("invalid number of teams")
	}

	teams := make(models.Teams, count)
	for i := 0; i < count; i++ {
		teams[i] = models.Team{
			Code:       helpers.NewCode(4),
			InstanceID: user.CurrentInstanceID,
		}
	}
	_, err := db.DB.NewInsert().Model(&teams).Exec(ctx)
	return err
}

func (s *GameManagerService) GetAllTeams(ctx context.Context, user *models.User) (models.Teams, error) {
	return models.FindAllTeams(ctx)
}

func (s *GameManagerService) GetAllLocations(ctx context.Context) (models.Locations, error) {
	return models.FindAllLocations(ctx)
}

func (s *GameManagerService) GetTeamActivityOverview(ctx context.Context) ([]map[string]interface{}, error) {
	return models.TeamActivityOverview(ctx)
}

func (s *GameManagerService) SaveLocation(ctx context.Context, location *models.Location, lat, lng, name string) error {
	if lat != "" && lng != "" {
		latFloat, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			return err
		}
		lngFloat, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			return err
		}
		location.Marker.SetCoords(latFloat, lngFloat)
		location.Marker.Save(ctx)
	}
	return location.Save(ctx)
}

func (s *GameManagerService) CreateLocation(ctx context.Context, user *models.User, name, content, criteriaID, lat, lng string) error {
	location := &models.Location{
		Name:       name,
		InstanceID: user.CurrentInstanceID,
		CriteriaID: criteriaID,
	}

	locationContent := models.LocationContent{
		Content: content,
	}

	if err := locationContent.Save(ctx); err != nil {
		return err
	}
	location.ContentID = locationContent.ID

	var latFloat, lngFloat float64
	var err error
	if lat != "" && lng != "" {
		latFloat, err = strconv.ParseFloat(lat, 64)
		if err != nil {
			return err
		}
		lngFloat, err = strconv.ParseFloat(lng, 64)
		if err != nil {
			return err
		}
	}

	marker := &models.Marker{
		Name: name,
		Lat:  latFloat,
		Lng:  lngFloat,
	}

	if err := marker.Save(ctx); err != nil {
		return err
	}
	location.MarkerID = marker.Code

	return location.Save(ctx)
}

func (s *GameManagerService) UpdateLocation(ctx context.Context, location *models.Location, newName, newContent, lat, lng string) error {
	if newContent != "" {
		location.Content.Content = newContent
		if err := location.Content.Save(ctx); err != nil {
			return err
		}
	}

	marker := location.Marker
	if lat != "" && lng != "" {
		latFloat, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			return err
		}
		lngFloat, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			return err
		}

		// Check if the marker is shared
		shared, err := s.isMarkerShared(ctx, location.MarkerID, location.InstanceID)
		if err != nil {
			return err
		}

		if shared {
			// Create a new marker since it's shared
			newMarker := &models.Marker{
				Name: marker.Name, // Keep the same name
				Lat:  latFloat,
				Lng:  lngFloat,
			}
			if err := newMarker.Save(ctx); err != nil {
				return err
			}
			location.MarkerID = newMarker.Code
		} else {
			// Update existing marker's coordinates and name
			marker.Lat = latFloat
			marker.Lng = lngFloat
			marker.Name = newName // Update the marker name if not shared
			if err := marker.Save(ctx); err != nil {
				return err
			}
		}
	}

	location.Name = newName
	return location.Save(ctx)
}

func (s *GameManagerService) isMarkerShared(ctx context.Context, markerID, instanceID string) (bool, error) {
	count, err := db.DB.NewSelect().
		Model((*models.Location)(nil)).
		Where("marker_id = ? AND instance_id != ?", markerID, instanceID).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
