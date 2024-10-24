package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/helpers"
	internalModels "github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

type GameManagerService struct {
	locationService LocationService
	userService     UserService
	teamService     TeamService
}

func NewGameManagerService() GameManagerService {
	return GameManagerService{
		locationService: NewLocationService(repositories.NewClueRepository()),
		userService:     NewUserService(repositories.NewUserRepository()),
		teamService:     NewTeamService(repositories.NewTeamRepository()),
	}
}

func (s *GameManagerService) CreateInstance(ctx context.Context, name string, user *internalModels.User) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	if name == "" {
		response.AddFlashMessage(*flash.NewError("Name is required"))
		response.Error = errors.New("name is required")
	}

	instance := &internalModels.Instance{
		Name:   name,
		UserID: user.ID,
	}

	if err := instance.Save(ctx); err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving instance"))
		response.Error = fmt.Errorf("saving instance: %w", err)
		return response
	}

	user.CurrentInstanceID = instance.ID
	err := s.userService.UpdateUser(ctx, user)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error updating user"))
		response.Error = fmt.Errorf("updating user: %w", err)
		return response
	}

	settings := &internalModels.InstanceSettings{
		InstanceID: instance.ID,
	}
	err = settings.Save(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving settings"))
		response.Error = fmt.Errorf("saving settings: %w", err)
		return response
	}

	response.Data["instanceID"] = instance.ID
	return response
}

func (s *GameManagerService) SwitchInstance(ctx context.Context, user *internalModels.User, instanceID string) (*internalModels.Instance, error) {
	instance, err := internalModels.FindInstanceByID(ctx, instanceID)
	if err != nil {
		return nil, errors.New("instance not found")
	}

	user.CurrentInstanceID = instance.ID
	if s.userService == nil {
		return nil, errors.New("user service is nil")
	}
	err = s.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return instance, nil
}

// Duplicate an instance
// This will create a new instance with the same name and locations
// The teams will not be duplicated
func (s *GameManagerService) DuplicateInstance(ctx context.Context, user *internalModels.User, id, name string) (response ServiceResponse) {
	response = ServiceResponse{}
	oldInstance, err := internalModels.FindInstanceByID(ctx, id)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Instance not found"))
		response.Error = fmt.Errorf("finding instance: %w", err)
		return response
	}

	locations, err := internalModels.FindAllLocations(ctx, oldInstance.ID)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error finding locations"))
		response.Error = fmt.Errorf("finding locations: %w", err)
		return response
	}

	newInstance := &internalModels.Instance{
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
		newLocation := &internalModels.Location{
			Name:       location.Name,
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

func (s *GameManagerService) LoadTeams(ctx context.Context, teams *[]internalModels.Team) error {
	for i := range *teams {
		err := s.teamService.LoadRelation(ctx, &(*teams)[i], "Scans")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GameManagerService) DeleteInstance(ctx context.Context, user *internalModels.User, instanceID, confirmName string) (response ServiceResponse) {
	response = ServiceResponse{}
	instance, err := internalModels.FindInstanceByID(ctx, instanceID)
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
		err := s.userService.UpdateUser(ctx, user)
		if err != nil {
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
	response.Data = make(map[string]interface{})
	if count < 1 {
		response.AddFlashMessage(*flash.NewError("Please enter a valid number of teams (1 or more)"))
		return response
	}

	teams := make([]internalModels.Team, count)
	for i := 0; i < count; i++ {
		teams[i] = internalModels.Team{
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
	response.Data["teams"] = teams
	return response
}

func (s *GameManagerService) GetAllLocations(ctx context.Context, instanceID string) ([]internalModels.Location, error) {
	return internalModels.FindAllLocations(ctx, instanceID)
}

func (s *GameManagerService) GetTeamActivityOverview(ctx context.Context, instanceID string) ([]map[string]interface{}, error) {
	return internalModels.TeamActivityOverview(ctx, instanceID)
}

func (s *GameManagerService) SaveLocation(ctx context.Context, location *internalModels.Location, lat, lng, name string) error {
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

func (s *GameManagerService) CreateLocation(ctx context.Context, user *internalModels.User, data map[string]string) (response *ServiceResponse) {

	name := data["name"]
	lat := data["latitude"]
	lng := data["longitude"]

	response = &ServiceResponse{}
	location := &internalModels.Location{
		Name:       name,
		InstanceID: user.CurrentInstanceID,
	}

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

	marker := &internalModels.Marker{
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

func (s *GameManagerService) UpdateLocation(ctx context.Context, location *internalModels.Location, newName, newContent, lat, lng string, points int) error {
	location.Points = points

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
			newMarker := &internalModels.Marker{
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
		Model((*internalModels.Location)(nil)).
		Where("marker_id = ? AND instance_id != ?", markerID, instanceID).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *GameManagerService) ValidateLocationMarker(user *internalModels.User, id string) bool {
	for _, loc := range user.CurrentInstance.Locations {
		if loc.MarkerID == id {
			return true
		}
	}
	return false
}

func (s *GameManagerService) ValidateLocationID(user *internalModels.User, id string) bool {
	for _, loc := range user.CurrentInstance.Locations {
		if loc.ID == id {
			return true
		}
	}
	return false
}

func (s *GameManagerService) GetQRCodePathAndContent(action, id, name, extension string) (string, string) {
	content := os.Getenv("SITE_URL")
	path := "assets/codes/"
	name = strings.Trim(name, " ")
	re := regexp.MustCompile(`[^\d\p{Latin} -]`)
	name = re.ReplaceAllString(name, "")
	if action == "in" {
		content = content + "/s/" + id
		path = path + extension + "/" + id + " " + name + "." + extension
	} else {
		content = content + "/o/" + id
		path = path + extension + "/" + id + " " + name + " Check Out." + extension
	}
	return path, content
}

// UpdateSettings parses the form values and updates the instance settings
func (s *GameManagerService) UpdateSettings(ctx context.Context, settings *internalModels.InstanceSettings, form url.Values) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	// Navigation mode
	navMode, err := internalModels.ParseNavigationMode(form.Get("navigationMode"))
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Something went wrong parsing navigation mode. Please try again."))
		response.Error = fmt.Errorf("parsing navigation mode: %w", err)
		return response
	}
	settings.NavigationMode = navMode

	// Completion method
	completionMethod, err := internalModels.ParseCompletionMethod(form.Get("completionMethod"))
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Something went wrong parsing completion method. Please try again."))
		response.Error = fmt.Errorf("parsing completion method: %w", err)
		return response
	}
	settings.CompletionMethod = completionMethod

	// Navigation method
	navMethod, err := internalModels.ParseNavigationMethod(form.Get("navigationMethod"))
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

	// Enable points
	enablePoints := form.Has("enablePoints")
	settings.EnablePoints = enablePoints

	// Enable Bonus Points
	enableBonusPoints := form.Has("enableBonusPoints")
	settings.EnableBonusPoints = enableBonusPoints

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

// StartGame starts the game immediately
func (s *GameManagerService) StartGame(ctx context.Context, user *internalModels.User) (response ServiceResponse) {
	return s.SetStartTime(ctx, user, time.Now())
}

// StopGame stops the game immediately
func (s *GameManagerService) StopGame(ctx context.Context, user *internalModels.User) (response ServiceResponse) {
	return s.SetEndTime(ctx, user, time.Now())
}

// SetStartTime sets the game start time to the given time
func (s *GameManagerService) SetStartTime(ctx context.Context, user *internalModels.User, time time.Time) (response ServiceResponse) {
	response = ServiceResponse{}

	// Check if the game is already active
	if user.CurrentInstance.GetStatus() == internalModels.Active {
		response.AddFlashMessage(*flash.NewInfo("Game is already active"))
		response.Error = errors.New("game is already active")
		return response
	}

	// Update the start time
	user.CurrentInstance.StartTime = bun.NullTime{Time: time}
	if !user.CurrentInstance.EndTime.After(time) {
		user.CurrentInstance.EndTime = bun.NullTime{}
	}

	err := user.CurrentInstance.Update(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error updating instance"))
		response.Error = fmt.Errorf("updating instance with new time: %w", err)
		return response
	}

	msg := flash.NewSuccess("Game scheduled to start at " + time.Format("2006-01-02 15:04:05"))
	response.AddFlashMessage(*msg)
	return response
}

// SetEndTime sets the game end time to the given time
func (s *GameManagerService) SetEndTime(ctx context.Context, user *internalModels.User, time time.Time) (response ServiceResponse) {
	response = ServiceResponse{}

	// Check if the game is already closed
	if user.CurrentInstance.GetStatus() == internalModels.Closed {
		response.AddFlashMessage(*flash.NewError("Game is already over"))
		response.Error = errors.New("game is already closed")
		return response
	}

	// Make sure the end time is after the start time
	if !user.CurrentInstance.StartTime.IsZero() && time.Before(user.CurrentInstance.StartTime.Time) {
		response.AddFlashMessage(*flash.NewError("End time must be after start time"))
		response.Error = errors.New("end time must be after start time")
		return response
	}

	// Update the end time
	user.CurrentInstance.EndTime = bun.NullTime{Time: time}
	err := user.CurrentInstance.Update(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error updating instance"))
		response.Error = fmt.Errorf("updating instance with new time: %w", err)
		return response
	}

	msg := flash.NewSuccess("Game scheduled to end at " + time.Format("2006-01-02 15:04:05"))
	response.AddFlashMessage(*msg)
	return response
}

// ScheduleGame schedules the game to start and/or end at a specific time
// Expects a form with set_start, utc_start_date, utc_start_time, set_end, utc_end_date, and utc_end_time
func (s *GameManagerService) ScheduleGame(ctx context.Context, user *internalModels.User, start time.Time, end time.Time) (response ServiceResponse) {
	response = ServiceResponse{}

	instance := user.CurrentInstance

	// Check if the end time is before the start time
	if !instance.EndTime.IsZero() && instance.EndTime.Time.Before(instance.StartTime.Time) {
		response.AddFlashMessage(*flash.NewError("End time must be after start time"))
		response.Error = errors.New("end time must be after start time")
		return response
	}

	instance.StartTime = bun.NullTime{Time: start}
	instance.EndTime = bun.NullTime{Time: end}

	// Save instance
	err := instance.Update(ctx)
	if err != nil {
		response.AddFlashMessage(*flash.NewError("Error saving instance"))
		response.Error = fmt.Errorf("saving instance: %w", err)
		return response
	}

	return response
}

// ReorderLocations takes a list of location IDs and updates the order
func (s *GameManagerService) ReorderLocations(ctx context.Context, user *internalModels.User, codes []string) (response ServiceResponse) {
	response = ServiceResponse{}

	// Loop through the locations and update the order
	for i, locationID := range codes {
		location, err := internalModels.FindLocationByInstanceAndCode(ctx, user.CurrentInstanceID, locationID)
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

// UpdateClues updates the clues for a location
// The clues are passed as a slice of strings and the IDs are passed as a slice of strings
// There may be new clues, updated clues, or deleted clues
func (s *GameManagerService) UpdateClues(ctx context.Context, location *internalModels.Location, clues []string, ids []string) error {
	clueRepo := repositories.NewClueRepository()

	err := s.locationService.LoadCluesForLocation(ctx, location)
	if err != nil {
		return fmt.Errorf("loading clues for location: %w", err)
	}

	// Loop through the clues and update them
	for i, clue := range clues {
		if i < len(ids) {
			// Delete any empty clues
			if clue == "" {
				err := clueRepo.Delete(ctx, &location.Clues[i])
				if err != nil {
					return fmt.Errorf("deleting clue: %w", err)
				}
				continue
			}

			// Update existing clue
			location.Clues[i].Content = clue
			err := clueRepo.Save(ctx, &location.Clues[i])
			if err != nil {
				return fmt.Errorf("saving clue: %w", err)
			}
			continue
		}

		// Skip empty clues
		if clue == "" {
			continue
		}

		// Create new clue
		newClue := &models.Clue{
			InstanceID: location.InstanceID,
			LocationID: location.ID,
			Content:    clue,
		}
		err := clueRepo.Save(ctx, newClue)
		if err != nil {
			return fmt.Errorf("saving new clue: %w", err)
		}
	}

	// Delete any remaining clues
	for i := len(clues); i < len(location.Clues); i++ {
		err := clueRepo.Delete(ctx, &location.Clues[i])
		if err != nil {
			return fmt.Errorf("deleting clue: %w", err)
		}
	}

	return nil
}

// DeleteLocation deletes a location
func (s *GameManagerService) DeleteLocation(ctx context.Context, location *internalModels.Location) error {
	return location.Delete(ctx)
}
