package models

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/ScanScout/helpers"
	"github.com/uptrace/bun"
)

type Team struct {
	baseModel

	InstanceID       string   `bun:",notnull" json:"instance_id"`
	Instance         Instance `bun:"rel:has-one,join:instance_id=id" json:"instance"`
	Code             string   `bun:",unique,pk" json:"code"`
	Scans            Scans    `bun:"rel:has-many,join:code=team_id" json:"scans"`
	MustScanOut      string   `bun:"" json:"must_scan_out"`
	BlockingLocation Location `bun:"rel:has-one,join:must_scan_out=code" json:"blocking_location"`
}

type Teams []Team

// FindAll returns all teams
func FindAllTeams(ctx context.Context) (Teams, error) {
	user := GetUserFromContext(ctx)
	var teams Teams
	err := db.NewSelect().
		Model(&teams).
		Where("team.instance_id = ?", user.CurrentInstanceID).
		Scan(ctx)
	return teams, err
}

// FindTeamByCode returns a team by code
func FindTeamByCode(ctx context.Context, code string) (*Team, error) {
	code = strings.ToUpper(code)
	var team Team
	err := db.NewSelect().Model(&team).Where("team.code = ? AND team.instance_id = ?", code).
		Relation("Scans").
		Relation("Scans.Location").
		Relation("BlockingLocation").
		Limit(1).Scan(ctx)
	return &team, err
}

// FindTeamByCodeAndInstance returns a team by code
func FindTeamByCodeAndInstance(ctx context.Context, code, instance string) (*Team, error) {
	code = strings.ToUpper(code)
	var team Team
	err := db.NewSelect().Model(&team).Where("team.code = ? and team.instance_id = ?", code, instance).
		Relation("BlockingLocation").
		Relation("Scans", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("instance_id = ?", instance).Order("name ASC")
		}).
		Limit(1).Scan(ctx)
	return &team, err
}

// HasVisited returns true if the team has visited the given location
func (t *Team) HasVisited(location *Location) bool {
	for _, s := range t.Scans {
		if s.LocationID == location.Code {
			return true
		}
	}
	return false
}

// SuggestNextLocation returns the next location to scan in
func (t *Team) SuggestNextLocations(ctx context.Context, limit int) *Locations {
	var locations Locations

	// Get the list of locations the team has already visited
	visited := make([]string, len(t.Scans))
	for i, s := range t.Scans {
		visited[i] = s.LocationID
	}

	var err error
	if len(visited) != 0 {
		// Get the list of locations the team must visit
		err = db.NewSelect().Model(&locations).
			Where("location.instance_id = ?", t.InstanceID).
			Where("location.code NOT IN (?)", bun.In(visited)).
			Scan(ctx)
	} else {
		err = db.NewSelect().Model(&locations).
			Where("location.instance_id = ?", t.InstanceID).
			Scan(ctx)
	}
	if err != nil {
		log.Error(err)
		return nil
	}

	seed := t.Code + fmt.Sprintf("%v", visited)
	h := fnv.New64a()
	_, err = h.Write([]byte(seed))
	if err != nil {
		log.Error(err)
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
func AddTeams(ctx context.Context, count int) error {
	user := GetUserFromContext(ctx)
	teams := make(Teams, count)
	for i := 0; i < count; i++ {
		teams[i] = Team{
			Code:       helpers.NewCode(4),
			InstanceID: user.CurrentInstanceID,
		}
	}
	_, err := db.NewInsert().Model(&teams).Exec(ctx)
	return err
}

func (t *Team) Update(ctx context.Context) error {
	_, err := db.NewUpdate().Model(t).WherePK().Exec(ctx)
	return err
}

// TeamActivityOverview returns a list of teams and their activity
func TeamActivityOverview(ctx context.Context) ([]map[string]interface{}, error) {
	// Get all instanceLocations
	instanceLocations, err := FindAllInstanceLocations(ctx)
	if err != nil {
		return nil, err
	}

	// Query all teams which have visited a location
	var teams Teams
	err = db.NewSelect().Model(&teams).
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
		if len(team.Scans) > 0 {
			count++
		}
	}
	activity := make([]map[string]interface{}, count)
	count = 0
	for _, team := range teams {
		if len(team.Scans) == 0 {
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
		for _, scan := range team.Scans {
			for j, instanceLocation := range instanceLocations {
				if scan.LocationID == instanceLocation.Location.Code {
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
