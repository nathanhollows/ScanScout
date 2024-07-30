package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type GameManagerService struct{}

func NewGameManagerService() *GameManagerService {
	return &GameManagerService{}
}

func (s *GameManagerService) CreateInstance(ctx context.Context, name string, user *models.User) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	if name == "" {
		response.AddFlashMessage(*flash.NewError("Name is required"))
		response.Error = errors.New("name is required")
	}

	instance := &models.Instance{
		Name:   name,
		UserID: user.ID,
	}

	if err := instance.Save(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving instance"))
		response.Error = fmt.Errorf("saving instance: %w", err)
		return response
	}

	user.CurrentInstanceID = instance.ID
	if err := user.Update(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error updating user"))
		response.Error = fmt.Errorf("updating user: %w", err)
		return response
	}

	settings := &models.InstanceSettings{
		InstanceID: instance.ID,
	}
	err := settings.Save(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving settings"))
		response.Error = fmt.Errorf("saving settings: %w", err)
		return response
	}

	response.Data["instanceID"] = instance.ID
	return response
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

// Duplicate an instance
// This will create a new instance with the same name and locations
// The teams will not be duplicated
func (s *GameManagerService) DuplicateInstance(ctx context.Context, user *models.User, id, name string) (response ServiceResponse) {
	response = ServiceResponse{}
	oldInstance, err := models.FindInstanceByID(ctx, id)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Instance not found"))
		response.Error = fmt.Errorf("finding instance: %w", err)
		return response
	}

	locations, err := models.FindAllLocations(ctx, oldInstance.ID)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error finding locations"))
		response.Error = fmt.Errorf("finding locations: %w", err)
		return response
	}

	newInstance := &models.Instance{
		Name:   name,
		UserID: user.ID,
	}

	if err := newInstance.Save(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving new instance"))
		response.Error = fmt.Errorf("saving new instance: %w", err)
		return response
	}

	// Copy locations
	for _, location := range locations {
		newContent := &models.LocationContent{
			Content: location.Content.Content,
		}
		if err := newContent.Save(ctx); err != nil {
			response.AddFlashMessage(*flash.NewError("Error saving location content: " + location.Name))
			response.Error = fmt.Errorf("saving location content: %w", err)
			return response
		}

		newLocation := &models.Location{
			Name:       location.Name,
			ContentID:  newContent.ID,
			InstanceID: newInstance.ID,
			MarkerID:   location.MarkerID,
		}
		if err := newLocation.Save(ctx); err != nil {
			response.AddFlashMessage(*flash.NewError("Error saving location: " + location.Name))
			response.Error = fmt.Errorf("saving location: %w", err)
			return response
		}
	}

	// Copy settings
	settings := oldInstance.Settings
	settings.InstanceID = newInstance.ID
	if err := settings.Save(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving settings"))
		response.Error = fmt.Errorf("saving settings: %w", err)
		return response
	}

	response.Data = make(map[string]interface{})
	response.Data["instanceID"] = newInstance.ID
	return response
}

func (s *GameManagerService) DeleteInstance(ctx context.Context, user *models.User, instanceID, confirmName string) (response ServiceResponse) {
	response = ServiceResponse{}
	instance, err := models.FindInstanceByID(ctx, instanceID)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Instance not found"))
		response.Error = fmt.Errorf("finding instance: %w", err)
		return response
	}

	if user.ID != instance.UserID {
		response.AddFlashMessage(*flash.NewError("You do not have permission to delete this instance"))
		response.Error = errors.New("you do not have permission to delete this instance")
		return response
	}

	if confirmName != instance.Name {
		response.AddFlashMessage(*flash.NewError("Instance name does not match confirmation"))
		response.Error = errors.New("instance name does not match confirmation")
		return response
	}

	if user.CurrentInstanceID == instance.ID {
		user.CurrentInstanceID = ""
		if err := user.Update(ctx); err != nil {
			response.AddFlashMessage(*flash.NewError("Error updating user"))
			response.Error = fmt.Errorf("updating user: %w", err)
			return response
		}
	}

	if err := instance.Delete(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error deleting instance"))
		response.Error = fmt.Errorf("deleting instance: %w", err)
		return response
	}

	response.AddFlashMessage(*flash.NewSuccess("Instance deleted!"))
	return response
}

func (s *GameManagerService) AddTeams(ctx context.Context, instanceID string, count int) (response ServiceResponse) {
	response = ServiceResponse{}
	if count < 1 {
		response.AddFlashMessage(*flash.NewError("Please enter a valid number of teams (1 or more)"))
		return response
	}

	teams := make(models.Teams, count)
	for i := 0; i < count; i++ {
		teams[i] = models.Team{
			Code:       helpers.NewCode(4),
			InstanceID: instanceID,
		}
	}
	_, err := db.DB.NewInsert().Model(&teams).Exec(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error adding teams"))
		response.Error = fmt.Errorf("AddTeams save teams: %w", err)
		return response
	}
	response.AddFlashMessage(*flash.NewSuccess("Teams added!"))
	return response
}

func (s *GameManagerService) GetAllLocations(ctx context.Context, instanceID string) (models.Locations, error) {
	return models.FindAllLocations(ctx, instanceID)
}

func (s *GameManagerService) GetTeamActivityOverview(ctx context.Context, instanceID string) ([]map[string]interface{}, error) {
	return models.TeamActivityOverview(ctx, instanceID)
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

func (s *GameManagerService) CreateLocation(ctx context.Context, user *models.User, name, content, criteriaID, lat, lng string) (response *ServiceResponse) {
	response = &ServiceResponse{}
	location := &models.Location{
		Name:       name,
		InstanceID: user.CurrentInstanceID,
	}

	locationContent := models.LocationContent{
		Content: content,
	}

	if err := locationContent.Save(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving location content: " + err.Error()))
		response.Error = err
		return response
	}
	location.ContentID = locationContent.ID

	var latFloat, lngFloat float64
	var err error
	if lat != "" && lng != "" {
		latFloat, err = strconv.ParseFloat(lat, 64)
		if err != nil {
			response.AddFlashMessage(*flash.NewError("Something went wrong parsing coordinates. Please try again."))
			response.Error = err
			return response
		}
		lngFloat, err = strconv.ParseFloat(lng, 64)
		if err != nil {
			response.AddFlashMessage(*flash.NewError("Something went wrong parsing coordinates. Please try again."))
			response.Error = err
			return response
		}
	}

	marker := &models.Marker{
		Name: name,
		Lat:  latFloat,
		Lng:  lngFloat,
	}

	if err := marker.Save(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving marker. Please try editing the location again."))
		response.Error = err
		return response
	}
	location.MarkerID = marker.Code
	location.Save(ctx)

	response.AddFlashMessage(*flash.NewSuccess("Location added!"))
	response.Data = make(map[string]interface{})
	response.Data["location"] = location
	return response
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

// UpdateSettings parses the form values and updates the instance settings
func (s *GameManagerService) UpdateSettings(ctx context.Context, settings *models.InstanceSettings, form url.Values) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	// Navigation mode
	navMode, err := models.ParseNavigationMode(form.Get("navigationMode"))
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Something went wrong parsing navigation mode. Please try again."))
		response.Error = fmt.Errorf("parsing navigation mode: %w", err)
		return response
	}
	settings.NavigationMode = navMode

	// Completion method
	completionMethod, err := models.ParseCompletionMethod(form.Get("completionMethod"))
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Something went wrong parsing completion method. Please try again."))
		response.Error = fmt.Errorf("parsing completion method: %w", err)
		return response
	}
	settings.CompletionMethod = completionMethod

	// Navigation method
	navMethod, err := models.ParseNavigationMethod(form.Get("navigationMethod"))
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Something went wrong parsing navigation method. Please try again."))
		response.Error = fmt.Errorf("parsing navigation method: %w", err)
		return response
	}
	settings.NavigationMethod = navMethod

	// Show team count
	showTeamCount := form.Has("showTeamCount")
	settings.ShowTeamCount = showTeamCount

	// Max locations
	maxLoc := form.Get("maxLocations")
	if maxLoc != "" {
		maxLocInt, err := strconv.Atoi(form.Get("maxLocations"))
		if err != nil {
			response.AddFlashMessage(*flash.NewError("Something went wrong parsing max locations. Please try again."))
			response.Error = fmt.Errorf("parsing max locations: %w", err)
			return response
		}
		settings.MaxNextLocations = maxLocInt
	}

	// Save settings
	err = settings.Save(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving settings. Please try again."))
		response.Error = fmt.Errorf("saving settings: %w", err)
		return
	}

	response.AddFlashMessage(*flash.NewSuccess("Settings updated!"))
	return response
}

// ReorderLocations takes a list of location IDs and updates the order
func (s *GameManagerService) ReorderLocations(ctx context.Context, user *models.User, codes []string) (response ServiceResponse) {
	response = ServiceResponse{}

	// Loop through the locations and update the order
	for i, locationID := range codes {
		location, err := models.FindLocationByInstanceAndCode(ctx, user.CurrentInstanceID, locationID)
		if err != nil {
			response.AddFlashMessage(*flash.NewError("Error finding location: " + locationID))
			response.Error = fmt.Errorf("finding location: %w", err)
			return response
		}
		location.Order = i
		if err := location.Save(ctx); err != nil {
			response.AddFlashMessage(*flash.NewError("Error saving location: " + locationID))
			response.Error = fmt.Errorf("saving location: %w", err)
			return response
		}
	}

	response.AddFlashMessage(*flash.NewSuccess("Locations reordered!"))
	return response
}
