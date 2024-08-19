package models

import (
	"context"

	"github.com/nathanhollows/Rapua/pkg/db"
)

type InstanceSettings struct {
	InstanceID        string           `bun:",pk,type:varchar(36)" json:"instance_id"`
	NavigationMode    NavigationMode   `bun:",type:int" json:"navigation_mode"`
	NavigationMethod  NavigationMethod `bun:",type:int" json:"navigation_method"`
	MaxNextLocations  int              `bun:",type:int,default:3" json:"max_next_locations"`
	CompletionMethod  CompletionMethod `bun:",type:int" json:"completion_method"`
	ShowTeamCount     bool             `bun:",type:bool" json:"show_team_count"`
	EnablePoints      bool             `bun:",type:bool" json:"enable_points"`
	EnableBonusPoints bool             `bun:",type:bool" json:"enable_bonus_points"`
}

// Save saves the instance settings to the database
func (s *InstanceSettings) Save(ctx context.Context) error {
	_, err := db.DB.NewInsert().Model(s).On("conflict (instance_id) do update").Exec(ctx)
	return err
}
