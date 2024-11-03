package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/services"
)

const (
	TierTrial    = "Trial"
	TierBasic    = "Basic"
	TierPro      = "Pro"
	TierEducator = "Educator"
)

type billingService struct {
	billingRepository repositories.BillingRepository
	userRepository    repositories.UserRepository
}

// NewBillingService creates a new billing service
func NewBillingService() services.BillingService {
	return &billingService{
		repositories.NewBillingRepository(),
		repositories.NewUserRepository(),
	}
}

// CheckPlanLimits checks the plan limits for a team
func (s *billingService) CheckPlanLimits(ctx context.Context, userID string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Event Boost gives users a temporary increase in their plan limits
	if time.Now().After(user.EventBoostExpiry.Time) {
		// Reset the event boost
		user.EventBoostExpiry = sql.NullTime{}
		s.userRepository.ResetEventBoost(ctx, userID)
	}

	// If the user has an event boost, they can play an unlimited number of games
	if user.EventBoostExpiry.Valid {
		return nil
	}

	// Even though the tiers are also limited by the number of players per game,
	// the billing service only checks the number of games played this month.
	// The number of players per game is checked by the game manager service.
	switch user.Tier {
	case TierTrial:
		if user.GamesThisMonth >= 3 {
			return fmt.Errorf("plan limit reached")
		}
	case TierBasic:
		if user.GamesThisMonth >= 10 {
			return fmt.Errorf("plan limit reached")
		}
	case TierPro, TierEducator:
		if user.GamesThisMonth >= 20 {
			return fmt.Errorf("plan limit reached")
		}
	}

	return nil
}

// SubscribeUser subscribes a user to a plan
func (s *billingService) SubscribeUser(ctx context.Context, userID string, tier string) error {
	_, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	return nil

}

// UpgradePlan upgrades the plan for a team
func (s *billingService) UpgradePlan(ctx context.Context, userID string, newTier string) error {
	return nil
}

// ActivateEventBoost activates an event boost for a team
func (s *billingService) ActivateEventBoost(ctx context.Context, userID string, duration time.Duration) error {
	return nil
}

// TrackUsage tracks the usage of a team's instance
func (s *billingService) TrackUsage(ctx context.Context, userID string, instanceID string, numPlayers int) error {
	return nil
}

// RequiresPlanSelection checks if a user needs to select a plan
func (s *billingService) RequiresPlanSelection(ctx context.Context, userID string) (bool, error) {
	status, err := s.GetPlanStatus(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("getting plan status: %w", err)
	}

	if status.Tier == "" {
		return true, nil
	}

	if (status.PlanExpiryDate.IsZero() || time.Now().After(status.PlanExpiryDate)) &&
		(status.EventBoostExpiry.IsZero() || time.Now().After(status.EventBoostExpiry)) {
		return true, nil
	}

	return false, nil
}

// GetPlanStatus retrieves the plan status for a user
func (s *billingService) GetPlanStatus(ctx context.Context, userID string) (*services.BillingStatus, error) {
	return s.billingRepository.GetPlanStatus(ctx, userID)
}
