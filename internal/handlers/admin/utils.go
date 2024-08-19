package handlers

import (
	"context"
	"log/slog"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

type AdminHandler struct {
	Logger              *slog.Logger
	GameManagerService  *services.GameManagerService
	NotificationService services.NotificationService
}

func NewAdminHandler(logger *slog.Logger, gs *services.GameManagerService, ns services.NotificationService) *AdminHandler {
	return &AdminHandler{
		Logger:              logger,
		GameManagerService:  gs,
		NotificationService: ns,
	}
}

// GetUserFromContext retrieves the user from the context
// User will always be in the context because the middleware
func (h AdminHandler) UserFromContext(ctx context.Context) *models.User {
	return ctx.Value(contextkeys.UserKey).(*models.User)
}
