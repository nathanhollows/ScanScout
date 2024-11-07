package repositories

import (
	"context"

	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/services"
)

type BillingRepository interface {
	// Retrieve billing information for the user
	GetPlanStatus(ctx context.Context, userID string) (*services.PlanStatus, error)
}

type billingRepository struct{}

// NewBillingRepository creates a new billing repository
func NewBillingRepository() BillingRepository {
	return &billingRepository{}
}

// GetBillingInfo retrieves the billing information for a user
func (r *billingRepository) GetPlanStatus(ctx context.Context, userID string) (*services.PlanStatus, error) {
	var billingStatus services.PlanStatus
	err := db.DB.NewSelect().
		Model(&models.User{}).
		Column("tier", "event_boost_expiry").
		Where("id = ?", userID).
		Scan(ctx, &billingStatus)
	if err != nil {
		return nil, err
	}

	return &billingStatus, nil
}
