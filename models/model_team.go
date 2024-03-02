package models

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/ScanScout/helpers"
	"github.com/uptrace/bun"
)

type Team struct {
	baseModel
	belongsToInstance

	Code             string   `bun:",unique,pk" json:"code"`
	Scans            Scans    `bun:"rel:has-many,join:code=team_id" json:"scans"`
	MustScanOut      string   `bun:"" json:"must_scan_out"`
	BlockingLocation Location `bun:"rel:has-one,join:must_scan_out=code" json:"blocking_location"`
}

type Teams []Team

// FindAll returns all teams
func FindAllTeams() (Teams, error) {
	ctx := context.Background()
	var teams Teams
	err := db.NewSelect().Model(&teams).Scan(ctx)
	return teams, err
}

// FindTeamByCode returns a team by code
func FindTeamByCode(code string) (*Team, error) {
	ctx := context.Background()
	var team Team
	err := db.NewSelect().Model(&team).Where("team.code = ?", code).
		Relation("Scans").
		Relation("Scans.Location").
		Relation("BlockingLocation").
		Limit(1).Scan(ctx)
	return &team, err
}

// FindTeamByCodeAndInstance returns a team by code
func FindTeamByCodeAndInstance(code, instance string) (*Team, error) {
	ctx := context.Background()
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
func (t *Team) SuggestNextLocations(limit int) *Locations {
	ctx := context.Background()
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
func AddTeams(count int) error {
	ctx := context.Background()
	teams := make(Teams, count)
	for i := 0; i < count; i++ {
		teams[i] = Team{
			Code: helpers.NewCode(4),
		}
	}
	_, err := db.NewInsert().Model(&teams).Exec(ctx)
	return err
}

func (t *Team) Update() error {
	ctx := context.Background()
	_, err := db.NewUpdate().Model(t).WherePK().Exec(ctx)
	return err
}
