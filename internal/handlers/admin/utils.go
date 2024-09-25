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
	UserServices        services.UserServices
	AssetGenerator      services.AssetGenerator
}

func NewAdminHandler(logger *slog.Logger, gs *services.GameManagerService, ns services.NotificationService, us services.UserServices) *AdminHandler {
	return &AdminHandler{
		Logger:              logger,
		GameManagerService:  gs,
		NotificationService: ns,
		UserServices:        us,
		AssetGenerator:      services.NewAssetGenerator(),
	}
}

// GetUserFromContext retrieves the user from the context
// User will always be in the context because the middleware
func (h AdminHandler) UserFromContext(ctx context.Context) *models.User {
	return ctx.Value(contextkeys.UserKey).(*models.User)
}
