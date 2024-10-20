package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/pkg/db"
)

type CheckInRepository interface {
	FindCheckInByTeamAndLocation(ctx context.Context, teamCode string, locationID string) (*models.Scan, error)
	Update(ctx context.Context, checkIn *models.Scan) error
	LogCheckIn(ctx context.Context, team models.Team, location models.Location, mustCheckOut bool, validationRequired bool) (models.Scan, error)
}

type checkInRepository struct{}

func NewCheckInRepository() CheckInRepository {
	return &checkInRepository{}
}

func (r *checkInRepository) FindCheckInByTeamAndLocation(ctx context.Context, teamCode string, locationID string) (*models.Scan, error) {
	var checkIn models.Scan
	err := db.DB.NewSelect().Model(&checkIn).Where("team_code = ? AND location_id = ?", teamCode, locationID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding check in by team and location: %w", err)
	}
	return &checkIn, nil
}

func (r *checkInRepository) Update(ctx context.Context, checkIn *models.Scan) error {
	_, err := db.DB.NewUpdate().Model(checkIn).WherePK().Exec(ctx)
	return err
}

// LogCheckIn logs a check in for a team at a location
func (r *checkInRepository) LogCheckIn(ctx context.Context, team models.Team, location models.Location, mustCheckOut bool, validationRequired bool) (models.Scan, error) {
	scan := &models.Scan{
		TeamID:          team.Code,
		LocationID:      location.ID,
		TimeIn:          time.Now().UTC(),
		MustScanOut:     mustCheckOut,
		Points:          location.Points,
		BlocksCompleted: !validationRequired,
	}
	err := scan.Save(ctx)
	if err != nil {
		return models.Scan{}, fmt.Errorf("saving scan: %w", err)
	}

	return *scan, nil
}
