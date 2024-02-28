package models

import (
	"context"

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
