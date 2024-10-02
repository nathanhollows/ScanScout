package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/services"
)

type PlayerHandler struct {
	Logger              *slog.Logger
	GameplayService     *services.GameplayService
	NotificationService services.NotificationService
	BlockService        services.BlockService
}

func NewPlayerHandler(logger *slog.Logger, gs *services.GameplayService, ns services.NotificationService) *PlayerHandler {
	return &PlayerHandler{
		Logger:              logger,
		GameplayService:     gs,
		NotificationService: ns,
		BlockService:        services.NewBlockService(repositories.NewBlockRepository()),
	}
}

// GetTeamFromContext retrieves the team from the context
// Team will always be in the context because the middleware
// However the Team could be nil if the team was not found
func (h PlayerHandler) getTeamFromContext(ctx context.Context) (*models.Team, error) {
	val := ctx.Value(contextkeys.TeamKey)
	if val == nil {
		return nil, errors.New("team not found")
	}
	team := val.(*models.Team)
	if team == nil {
		return nil, errors.New("team not found")
	}
	return team, nil
}

// redirect is a helper function to redirect the user to a new page
// It accounts for htmx requests and redirects the user to the referer
func (h PlayerHandler) redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		return
	}
	http.Redirect(w, r, path, http.StatusFound)
}
