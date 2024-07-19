package handlers

import (
	"context"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

type AdminHandler struct {
	GameManagerService *services.GameManagerService
}

func NewAdminHandler(gs *services.GameManagerService) *AdminHandler {
	return &AdminHandler{
		GameManagerService: gs,
	}
}

// GetUserFromContext retrieves the user from the context
// User will always be in the context because the middleware
func (h AdminHandler) UserFromContext(ctx context.Context) *models.User {
	return ctx.Value(contextkeys.UserIDKey).(*models.User)
}
