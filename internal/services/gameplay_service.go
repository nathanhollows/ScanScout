package services

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
	"github.com/uptrace/bun"
	"golang.org/x/exp/rand"
)

type GameplayService struct{}

func (s *GameplayService) GetTeamStatus(ctx context.Context, teamCode string) (*models.Team, error) {
	team, err := models.FindTeamByCode(ctx, teamCode)
	if err != nil {
		return nil, fmt.Errorf("GetTeamStatus: %w", err)
	}
	return team, nil
}

func (s *GameplayService) StartPlaying(ctx context.Context, teamCode, customTeamName string) (*models.Team, error) {
	team, err := models.FindTeamByCode(ctx, teamCode)
	if err != nil {
		return nil, fmt.Errorf("StartPlaying get team: %w", err)
	}

	// Update team with custom name if provided
	if !team.HasStarted || customTeamName != "" {
		team.Name = customTeamName
		team.HasStarted = true
		if err := team.Update(ctx); err != nil {
			return nil, fmt.Errorf("StartPlaying update team: %w", err)
		}
	}

	return team, nil
}

func (s *GameplayService) SuggestNextLocations(ctx context.Context, team *models.Team, limit int) ([]*models.Location, error) {
	var locations []*models.Location

	visited := make([]string, len(team.Scans))
	for i, s := range team.Scans {
		visited[i] = s.LocationID
	}

	var err error
	if len(visited) != 0 {
		err = db.DB.NewSelect().Model(&locations).
			Where("location.instance_id = ?", team.InstanceID).
			Where("location.code NOT IN (?)", bun.In(visited)).
			Relation("Marker").
			Scan(ctx)
	} else {
		err = db.DB.NewSelect().Model(&locations).
			Where("location.instance_id = ?", team.InstanceID).
			Relation("Marker").
			Scan(ctx)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}

	seed := team.Code + fmt.Sprintf("%s", visited)
	h := fnv.New64a()
	_, err = h.Write([]byte(seed))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	rand.New(rand.NewSource(uint64(h.Sum64()))).Shuffle(len(locations), func(i, j int) {
		locations[i], locations[j] = locations[j], locations[i]
	})

	if len(locations) > limit {
		locations = locations[:limit]
	}

	for i := 0; i < len(locations); i++ {
		for j := i + 1; j < len(locations); j++ {
			if locations[i].CurrentCount > locations[j].CurrentCount {
				locations[i], locations[j] = locations[j], locations[i]
			}
		}
	}

	return locations, nil
}
