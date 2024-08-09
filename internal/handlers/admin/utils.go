package handlers

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

type AdminHandler struct {
	GameManagerService  *services.GameManagerService
	NotificationService services.NotificationService
}

func NewAdminHandler(gs *services.GameManagerService, ns services.NotificationService) *AdminHandler {
	return &AdminHandler{
		GameManagerService:  gs,
		NotificationService: ns,
	}
}

// GetUserFromContext retrieves the user from the context
// User will always be in the context because the middleware
func (h AdminHandler) UserFromContext(ctx context.Context) *models.User {
	return ctx.Value(contextkeys.UserIDKey).(*models.User)
}
