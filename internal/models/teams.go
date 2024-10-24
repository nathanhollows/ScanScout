package models

import (
	"context"
	"fmt"

	"github.com/nathanhollows/Rapua/db"
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
