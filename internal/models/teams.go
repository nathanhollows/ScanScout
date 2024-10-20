package models

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/models"
	"github.com/uptrace/bun"
)

type Team struct {
	baseModel

	Code         string `bun:"code,unique,pk"`
	Name         string `bun:"name,"`
	InstanceID   string `bun:"instance_id,notnull"`
	HasStarted   bool   `bun:"has_started,default:false"`
	MustCheckOut string `bun:"must_scan_out"`
	Points       int    `bun:"points,"`

	Instance         Instance                `bun:"rel:has-one,join:instance_id=id"`
	CheckIns         []CheckIn               `bun:"rel:has-many,join:code=team_code"`
	BlockingLocation Location                `bun:"rel:has-one,join:must_scan_out=marker_id,join:instance_id=instance_id"`
	Messages         []Notification          `bun:"rel:has-many,join:code=team_code"`
	Blocks           []models.TeamBlockState `bun:"rel:has-many,join:code=team_code"`
}

// Delete removes the team from the database
func (t *Team) Delete(ctx context.Context) error {
	_, err := db.DB.NewDelete().Model(t).WherePK().Exec(ctx)
	return fmt.Errorf("Delete: %v", err)
}

// Updates the team into the database
func (t *Team) Update(ctx context.Context) error {
	_, err := db.DB.NewUpdate().Model(t).WherePK().Exec(ctx)
	return err
}

// FindAll returns all teams
func FindAllTeams(ctx context.Context, instanceID string) ([]Team, error) {
	var teams []Team
	err := db.DB.NewSelect().
		Model(&teams).
		Where("team.instance_id = ?", instanceID).
		// Add the scans in the relation order by location_id
		Relation("Scans", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("location_id ASC")
		}).
		Scan(ctx)
	if err != nil {
		return teams, fmt.Errorf("FindAllTeams: %w", err)
	}
	return teams, nil
}

// FindTeamByCode returns a team by code
func FindTeamByCode(ctx context.Context, code string) (*Team, error) {
	code = strings.ToUpper(code)
	var team Team
	err := db.DB.NewSelect().Model(&team).Where("team.code = ?", code).
		Relation("BlockingLocation").
		Relation("Instance").
		Relation("Instance.Settings").
		Limit(1).Scan(ctx)
	if err != nil {
		return &team, fmt.Errorf("FindTeamByCode: %v", err)
	}
	return &team, nil
}

// FindTeamByCodeAndInstance returns a team by code
func FindTeamByCodeAndInstance(ctx context.Context, code, instance string) (*Team, error) {
	code = strings.ToUpper(code)
	var team Team
	err := db.DB.NewSelect().Model(&team).Where("team.code = ? and team.instance_id = ?", code, instance).
		Relation("BlockingLocation").
		Relation("Scans", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("instance_id = ?", instance).Order("name ASC")
		}).
		Limit(1).Scan(ctx)
	return &team, fmt.Errorf("FindTeamByCodeAndInstance: %v", err)
}

// HasVisited returns true if the team has visited the given location
func (t *Team) HasVisited(location *Location) bool {
	for _, s := range t.CheckIns {
		if s.LocationID == location.ID {
			return true
		}
	}
	return false
}

// SuggestNextLocation returns the next location to scan in
func (t *Team) SuggestNextLocations(ctx context.Context, limit int) *Locations {
	var locations Locations

	// Get the list of locations the team has already visited
	visited := make([]string, len(t.CheckIns))
	for i, s := range t.CheckIns {
		visited[i] = s.LocationID
	}

	var err error
	if len(visited) != 0 {
		// Get the list of locations the team must visit
		err = db.DB.NewSelect().Model(&locations).
			Where("location.instance_id = ?", t.InstanceID).
			Where("location.code NOT IN (?)", bun.In(visited)).
			Scan(ctx)
	} else {
		err = db.DB.NewSelect().Model(&locations).
			Where("location.instance_id = ?", t.InstanceID).
			Scan(ctx)
	}
	if err != nil {
		return nil
	}

	seed := t.Code + fmt.Sprintf("%v", visited)
	h := fnv.New64a()
	_, err = h.Write([]byte(seed))
	if err != nil {
		return nil
	}

	rand.New(rand.NewSource(int64(h.Sum64()))).Shuffle(len(locations), func(i, j int) {
		locations[i], locations[j] = locations[j], locations[i]
	})

	// Limit the number of locations
	if len(locations) > limit {
		locations = locations[:limit]
	}

	// Order the locations by CurrentCount
	for i := 0; i < len(locations); i++ {
		for j := i + 1; j < len(locations); j++ {
			if locations[i].CurrentCount > locations[j].CurrentCount {
				locations[i], locations[j] = locations[j], locations[i]
			}
		}
	}

	return &locations
}

// AddTeams adds the given number of teams
func AddTeams(ctx context.Context, instanceID string, count int) error {
	teams := make([]Team, count)
	for i := 0; i < count; i++ {
		teams[i] = Team{
			Code:       helpers.NewCode(4),
			InstanceID: instanceID,
		}
	}
	_, err := db.DB.NewInsert().Model(&teams).Exec(ctx)
	return err
}

// TeamActivityOverview returns a list of teams and their activity
func TeamActivityOverview(ctx context.Context, instanceID string) ([]map[string]interface{}, error) {
	// Get all instanceLocations
	instanceLocations, err := FindAllLocations(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	// Query all teams which have visited a location
	var teams []Team
	err = db.DB.NewSelect().Model(&teams).
		Relation("Scans").
		Order("team.code ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// Construct an array of team activity
	// Each team has a list of every location, and a duration of time spent at each location
	// For each location we also mark if the team has visited it, is currently visiting it, or has not visited it
	// The duration is calculated by the time between TimeIn and TimeOut
	// If TimeOut is not set, the duration is calculated by the time between TimeIn and now
	// If TimeIn is not set, the duration is 0
	count := 0
	for _, team := range teams {
		if len(team.CheckIns) > 0 {
			count++
		}
	}
	activity := make([]map[string]interface{}, count)
	count = 0
	for _, team := range teams {
		if !team.HasStarted {
			continue
		}
		activity[count] = make(map[string]interface{})
		activity[count]["team"] = team
		activity[count]["locations"] = make([]map[string]interface{}, len(instanceLocations))
		for j, location := range instanceLocations {
			activity[count]["locations"].([]map[string]interface{})[j] = make(map[string]interface{})
			activity[count]["locations"].([]map[string]interface{})[j]["location"] = location
			activity[count]["locations"].([]map[string]interface{})[j]["visited"] = false
			activity[count]["locations"].([]map[string]interface{})[j]["visiting"] = false
			activity[count]["locations"].([]map[string]interface{})[j]["duration"] = 0
			activity[count]["locations"].([]map[string]interface{})[j]["time_in"] = ""
			activity[count]["locations"].([]map[string]interface{})[j]["time_out"] = ""
		}
		for _, scan := range team.CheckIns {
			for j, instanceLocation := range instanceLocations {
				if scan.LocationID == instanceLocation.Marker.Code {
					activity[count]["locations"].([]map[string]interface{})[j]["visited"] = true
					activity[count]["locations"].([]map[string]interface{})[j]["time_in"] = scan.TimeIn
					if scan.TimeOut.IsZero() {
						activity[count]["locations"].([]map[string]interface{})[j]["visiting"] = true
					} else {
						activity[count]["locations"].([]map[string]interface{})[j]["time_out"] = scan.TimeOut
						activity[count]["locations"].([]map[string]interface{})[j]["duration"] = scan.TimeOut.Sub(scan.TimeIn).Seconds()
					}
				}
			}
		}
		count++
	}

	return activity, err
}

// GetVisitedLocations returns locations visited by the team
func (t *Team) GetVisitedLocations(ctx context.Context) ([]*Location, error) {
	var locations []*Location
	err := db.DB.NewSelect().Model(&locations).
		Where("marker_id in (SELECT location_id FROM scans WHERE team_code = ?)", t.Code).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return locations, nil
}
