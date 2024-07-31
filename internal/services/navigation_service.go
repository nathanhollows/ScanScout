package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"golang.org/x/exp/rand"
)

type NavigationService struct{}

func NewNavigationService() *NavigationService {
	return &NavigationService{}
}

// CheckValidLocation checks if the location code is valid for the team to scan in to
// This function returns an error if the location code is invalid
func (s *NavigationService) CheckValidLocation(ctx context.Context, team *models.Team, settings *models.InstanceSettings, locationCode string) (response ServiceResponse) {
	response = ServiceResponse{}

	// Find valid locations
	resp := s.DetermineNextLocations(ctx, team)
	if resp.Error != nil {
		response.Error = fmt.Errorf("determine next valid locations: %w", resp.Error)
		return response
	}

	// Check if the location code is valid
	locationCode = strings.TrimSpace(strings.ToUpper(locationCode))
	for _, loc := range resp.Data["nextLocations"].(models.Locations) {
		if loc.Marker.Code == locationCode {
			return response
		}
	}

	response.Error = errors.New("invalid location")
	return response
}

func (s *NavigationService) DetermineNextLocations(ctx context.Context, team *models.Team) ServiceResponse {
	response := ServiceResponse{}

	err := team.LoadScans(ctx)
	if err != nil {
		response.Error = fmt.Errorf("getOrderedLocations load scans: %w", err)
		return response
	}

	err = team.LoadInstance(ctx)
	if err != nil {
		response.Error = fmt.Errorf("getOrderedLocations load instance: %w", err)
		return response
	}

	err = team.Instance.LoadLocations(ctx)
	if err != nil {
		response.Error = fmt.Errorf("getOrderedLocations load locations: %w", err)
		return response
	}

	// Check if the team has a blocking location
	if team.MustScanOut != "" {
		team.LoadBlockingLocation(ctx)
		response.Data["blockingLocation"] = team.BlockingLocation
		response.AddFlashMessage(*flash.NewInfo("You must scan out of " + team.BlockingLocation.Name + " before you can scan in to your next location."))
		return response
	}

	// Check if the team has visited all locations
	if len(team.Scans) == len(team.Instance.Locations) {
		response.Error = errors.New("all locations visited")
		response.AddFlashMessage(*flash.NewInfo("You have visited all locations!"))
		return response
	}

	// Determine the next locations based on the navigation mode
	switch team.Instance.Settings.NavigationMode {
	case models.OrderedNav:
		return s.getOrderedLocations(ctx, team)
	case models.RandomNav:
		return s.getRandomLocations(ctx, team)
	case models.FreeRoamNav:
		return s.getFreeRoamLocations(ctx, team)
	}

	response.Error = errors.New("invalid navigation mode")
	return response
}

// getUnvisitedLocations returns a list of locations that the team has not visited
func (s *NavigationService) getUnvisitedLocations(_ context.Context, team *models.Team) models.Locations {
	var unvisited models.Locations

	// Find the next location
	for _, location := range team.Instance.Locations {
		found := false
		for _, scan := range team.Scans {
			if scan.LocationID == location.Marker.Code {
				found = true
				break
			}
		}
		if !found {
			unvisited = append(unvisited, location)
			break
		}
	}

	return unvisited
}

// getOrderedLocations returns locations in the order defined by the admin
// This function returns the next location for the team to visit
func (s *NavigationService) getOrderedLocations(ctx context.Context, team *models.Team) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	unvisited := s.getUnvisitedLocations(ctx, team)
	if len(unvisited) == 0 {
		response.Error = errors.New("all locations visited")
		response.AddFlashMessage(*flash.NewInfo("You have visited all locations!"))
		return response
	}

	response.Data["nextLocations"] = unvisited[:1]
	fmt.Println(response.Data["nextLocations"])

	return response
}

// getRandomLocations returns random locations for the team to visit
// This function uses the team code as a seed for the random number generator
// Process:
// 1. Shuffle the list of all locations deterministically based on team code
// 2. Select the first n unvisited locations from the shuffled list
// 3. Return these locations ensuring the order is consistent across refreshes
func (s *NavigationService) getRandomLocations(ctx context.Context, team *models.Team) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	allLocations := team.Instance.Locations
	if len(allLocations) == 0 {
		response.Error = errors.New("no locations found")
		response.AddFlashMessage(*flash.NewInfo("No locations available"))
		return response
	}

	unvisited := s.getUnvisitedLocations(ctx, team)
	if len(unvisited) == 0 {
		response.Error = errors.New("all locations visited")
		response.AddFlashMessage(*flash.NewInfo("You have visited all locations!"))
		return response
	}

	// Seed the random number generator with the team code to ensure deterministic shuffling
	seed := uint64(0)
	for _, char := range team.Code {
		seed += uint64(char)
	}
	rand.Seed(seed)

	// Shuffle the list of all locations deterministically
	shuffledLocations := make([]models.Location, len(allLocations))
	copy(shuffledLocations, allLocations)
	rand.Shuffle(len(shuffledLocations), func(i, j int) {
		shuffledLocations[i], shuffledLocations[j] = shuffledLocations[j], shuffledLocations[i]
	})

	// Select the first n unvisited locations from the shuffled list
	n := team.Instance.Settings.MaxNextLocations
	var selectedLocations models.Locations
	for _, loc := range shuffledLocations {
		if !team.HasVisited(&loc) {
			selectedLocations = append(selectedLocations, loc)
			if len(selectedLocations) >= n {
				break
			}
		}
	}

	if len(selectedLocations) == 0 {
		response.Error = errors.New("no unvisited locations found")
		response.AddFlashMessage(*flash.NewInfo("No unvisited locations available"))
		return response
	}

	response.Data["nextLocations"] = selectedLocations
	return response
}

// getFreeRoamLocations returns a list of locations for free roam mode
// This function returns all locations in the instance for the team to visit
func (s *NavigationService) getFreeRoamLocations(ctx context.Context, team *models.Team) (response ServiceResponse) {
	response = ServiceResponse{}
	response.Data = make(map[string]interface{})

	unvisited := s.getUnvisitedLocations(ctx, team)

	if len(unvisited) == 0 {
		response.Error = errors.New("all locations visited")
		response.AddFlashMessage(*flash.NewInfo("You have visited all locations!"))
		return response
	}

	response.Data["nextLocations"] = unvisited
	return response
}
