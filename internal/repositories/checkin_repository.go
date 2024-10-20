package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/models"
)

type CheckInRepository interface {
	FindCheckInByTeamAndLocation(ctx context.Context, teamCode string, locationID string) (*models.CheckIn, error)
	Update(ctx context.Context, checkIn *models.CheckIn) error
	LogCheckIn(ctx context.Context, team models.Team, location models.Location, mustCheckOut bool, validationRequired bool) (models.CheckIn, error)
	CheckOut(ctx context.Context, team *models.Team, location *models.Location) (models.CheckIn, error)
}

type checkInRepository struct{}

func NewCheckInRepository() CheckInRepository {
	return &checkInRepository{}
}

func (r *checkInRepository) FindCheckInByTeamAndLocation(ctx context.Context, teamCode string, locationID string) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	err := db.DB.NewSelect().Model(&checkIn).Where("team_code = ? AND location_id = ?", teamCode, locationID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding check in by team and location: %w", err)
	}
	return &checkIn, nil
}

func (r *checkInRepository) Update(ctx context.Context, checkIn *models.CheckIn) error {
	_, err := db.DB.NewUpdate().Model(checkIn).WherePK().Exec(ctx)
	return err
}

// LogCheckIn logs a check in for a team at a location
func (r *checkInRepository) LogCheckIn(ctx context.Context, team models.Team, location models.Location, mustCheckOut bool, validationRequired bool) (models.CheckIn, error) {
	scan := &models.CheckIn{
		TeamID:          team.Code,
		LocationID:      location.ID,
		TimeIn:          time.Now().UTC(),
		MustCheckOut:    mustCheckOut,
		Points:          location.Points,
		BlocksCompleted: !validationRequired,
	}
	err := scan.Save(ctx)
	if err != nil {
		return models.CheckIn{}, fmt.Errorf("saving scan: %w", err)
	}

	return *scan, nil
}

// CheckOut logs a check out for a team at a location
func (r *checkInRepository) CheckOut(ctx context.Context, team *models.Team, location *models.Location) (models.CheckIn, error) {
	if team == nil {
		return models.CheckIn{}, fmt.Errorf("team is required")
	}

	if location == nil {
		return models.CheckIn{}, fmt.Errorf("location is required")
	}

	if len(team.CheckIns) == 0 {
		return models.CheckIn{}, fmt.Errorf("no check ins found for team")
	}

	var checkIn *models.CheckIn
	for i := range team.CheckIns {
		if team.CheckIns[i].LocationID == location.ID {
			checkIn = &team.CheckIns[i]
			break
		}
	}

	if checkIn == nil {
		return models.CheckIn{}, fmt.Errorf("check in not found")
	}

	checkIn.TimeOut = time.Now().UTC()
	checkIn.MustCheckOut = false
	err := r.Update(ctx, checkIn)
	if err != nil {
		return models.CheckIn{}, fmt.Errorf("updating check in: %w", err)
	}

	return *checkIn, nil
}
