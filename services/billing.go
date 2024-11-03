package services

import (
	"context"
	"time"
)

type BillingStatus struct {
	Tier             string
	PlanExpiryDate   time.Time
	EventBoostExpiry time.Time
}

type BillingService interface {
	// CheckPlanLimits checks if the user has exceeded their plan limits
	// and returns an error if they have
	CheckPlanLimits(ctx context.Context, userID string) error

	// UpgradePlan upgrades the user's plan to the new tier
	// and returns an error if the upgrade fails
	SubscribeUser(ctx context.Context, userID string, tier string) error

	// ActivateEventBoost activates an event boost for the user
	// and returns an error if the activation fails
	UpgradePlan(ctx context.Context, userID string, newTier string) error

	// SubscribeUser subscribes the user to the given tier
	// and returns an error if the subscription fails
	ActivateEventBoost(ctx context.Context, userID string, duration time.Duration) error

	// TrackUsage tracks the usage of the user's instance
	// and returns an error if the tracking fails
	TrackUsage(ctx context.Context, userID string, instanceID string, numPlayers int) error

	// Check if the user needs to select a plan (e.g. if their tier is empty)
	RequiresPlanSelection(ctx context.Context, userID string) (bool, error)

	// Retrieve the user's current plan status
	GetPlanStatus(ctx context.Context, userID string) (*BillingStatus, error)
}

type NoOpBillingService struct{}

// NewNoOpBillingService returns a new no-op billing service
func NewNoOpBillingService() BillingService {
	return &NoOpBillingService{}
}

func (n *NoOpBillingService) CheckPlanLimits(ctx context.Context, userID string) error {
	return nil
}

func (n *NoOpBillingService) SubscribeUser(ctx context.Context, userID string, tier string) error {
	return nil
}

func (n *NoOpBillingService) UpgradePlan(ctx context.Context, userID string, newTier string) error {
	return nil
}

func (n *NoOpBillingService) ActivateEventBoost(ctx context.Context, userID string, duration time.Duration) error {
	return nil
}

func (n *NoOpBillingService) TrackUsage(ctx context.Context, userID string, instanceID string, numPlayers int) error {
	return nil
}

func (n *NoOpBillingService) RequiresPlanSelection(ctx context.Context, userID string) (bool, error) {
	return false, nil
}

func (n *NoOpBillingService) GetPlanStatus(ctx context.Context, userID string) (*BillingStatus, error) {
	return &BillingStatus{
		Tier:             "",
		PlanExpiryDate:   time.Time{},
		EventBoostExpiry: time.Time{},
	}, nil
}
