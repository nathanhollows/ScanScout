package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nathanhollows/Rapua/internal/models"
	"golang.org/x/exp/rand"
)

var (
	ErrorAllLocationsVisited = errors.New("all locations visited")
	ErrInstanceNotFound      = errors.New("instance not found")
)

type NavigationService struct{}

func NewNavigationService() NavigationService {
	return NavigationService{}
}

// CheckValidLocation checks if the location code is valid for the team to scan in to
// This function returns an error if the location code is invalid
func (s NavigationService) CheckValidLocation(ctx context.Context, team *models.Team, settings *models.InstanceSettings, locationCode string) (bool, error) {

	// Find valid locations
	locations, err := s.DetermineNextLocations(ctx, team)
	if err != nil {
		return false, fmt.Errorf("determine next valid locations: %w", err)
	}

	// Check if the location code is valid
	locationCode = strings.TrimSpace(strings.ToUpper(locationCode))
	for _, loc := range locations {
		if loc.Marker.Code == locationCode {
			return true, nil
		}
	}
	return false, nil
}

func (s NavigationService) DetermineNextLocations(ctx context.Context, team *models.Team) ([]models.Location, error) {
	// Check if the team has visited all locations
	if len(team.CheckIns) == len(team.Instance.Locations) {
		return nil, ErrorAllLocationsVisited
	}

	if team.Instance.ID == "" || team.Instance.Settings.InstanceID == "" {
		return nil, ErrInstanceNotFound
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

	return nil, errors.New("invalid navigation mode")
}

// getUnvisitedLocations returns a list of locations that the team has not visited
func (s NavigationService) getUnvisitedLocations(_ context.Context, team *models.Team) []models.Location {
	unvisited := []models.Location{}

	// Find the next location
	for _, location := range team.Instance.Locations {
		found := false
		for _, scan := range team.CheckIns {
			if scan.LocationID == location.ID {
				found = true
				continue
			}
		}
		if !found {
			unvisited = append(unvisited, location)
			continue
		}
	}

	return unvisited
}

// getOrderedLocations returns locations in the order defined by the admin
// This function returns the next location for the team to visit
func (s NavigationService) getOrderedLocations(ctx context.Context, team *models.Team) ([]models.Location, error) {
	unvisited := s.getUnvisitedLocations(ctx, team)
	if len(unvisited) == 0 {
		return nil, ErrorAllLocationsVisited
	}

	// Reorder based on .Order
	for i := 0; i < len(unvisited); i++ {
		for j := i + 1; j < len(unvisited); j++ {
			if unvisited[i].Order > unvisited[j].Order {
				unvisited[i], unvisited[j] = unvisited[j], unvisited[i]
			}
		}
	}

	return unvisited[:1], nil
}

// getRandomLocations returns random locations for the team to visit
// This function uses the team code as a seed for the random number generator
// Process:
// 1. Shuffle the list of all locations deterministically based on team code
// 2. Select the first n unvisited locations from the shuffled list
// 3. Return these locations ensuring the order is consistent across refreshes
func (s NavigationService) getRandomLocations(ctx context.Context, team *models.Team) ([]models.Location, error) {
	allLocations := team.Instance.Locations
	if len(allLocations) == 0 {
		return nil, errors.New("no locations found")
	}

	unvisited := s.getUnvisitedLocations(ctx, team)
	if len(unvisited) == 0 {
		return nil, fmt.Errorf("all locations visited")
	}

	// Seed the random number generator with the team code to ensure deterministic shuffling
	seed := uint64(0)
	for _, char := range team.Code {
		seed += uint64(char)
	}
	rand.Seed(seed)

	// We shuffle the list of all locations to ensure randomness
	// even when the team has visited some locations
	shuffledLocations := make([]models.Location, len(allLocations))
	copy(shuffledLocations, allLocations)
	rand.Shuffle(len(shuffledLocations), func(i, j int) {
		shuffledLocations[i], shuffledLocations[j] = shuffledLocations[j], shuffledLocations[i]
	})

	// Select the first n unvisited locations from the shuffled list
	n := team.Instance.Settings.MaxNextLocations
	selectedLocations := []models.Location{}
	for _, loc := range shuffledLocations {
		if !s.HasVisited(team.CheckIns, loc.ID) {
			selectedLocations = append(selectedLocations, loc)
			if len(selectedLocations) >= n {
				break
			}
		}
	}

	if len(selectedLocations) == 0 {
		return nil, ErrorAllLocationsVisited
	}

	return selectedLocations, nil
}

// getFreeRoamLocations returns a list of locations for free roam mode
// This function returns all locations in the instance for the team to visit
func (s NavigationService) getFreeRoamLocations(ctx context.Context, team *models.Team) ([]models.Location, error) {
	unvisited := s.getUnvisitedLocations(ctx, team)

	if len(unvisited) == 0 {
		return nil, ErrorAllLocationsVisited
	}

	return unvisited, nil
}

// HasVisited returns true if the team has visited the location
func (s NavigationService) HasVisited(checkins []models.CheckIn, locationID string) bool {
	for _, checkin := range checkins {
		if checkin.LocationID == locationID {
			return true
		}
	}
	return false
}
