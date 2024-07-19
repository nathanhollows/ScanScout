package handlers

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

type PlayerHandler struct {
	GameplayService *services.GameplayService
}

func NewPlayerHandler(gs *services.GameplayService) *PlayerHandler {
	return &PlayerHandler{
		GameplayService: gs,
	}
}

// GetTeamFromContext retrieves the team from the context
// Team will always be in the context because the middleware
func (h PlayerHandler) getTeamFromContext(ctx context.Context) *models.Team {
	return ctx.Value(contextkeys.TeamKey).(*models.Team)
}
