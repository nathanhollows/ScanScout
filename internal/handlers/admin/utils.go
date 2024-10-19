package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

type AdminHandler struct {
	Logger              *slog.Logger
	GameManagerService  services.GameManagerService
	NotificationService services.NotificationService
	UserServices        services.UserServices
	AssetGenerator      services.AssetGenerator
	LocationService     services.LocationService
	BlockService        services.BlockService
}

func NewAdminHandler(logger *slog.Logger, gs services.GameManagerService, ns services.NotificationService, us services.UserServices) *AdminHandler {
	return &AdminHandler{
		Logger:              logger,
		GameManagerService:  gs,
		NotificationService: ns,
		UserServices:        us,
		AssetGenerator:      services.NewAssetGenerator(),
		LocationService:     services.NewLocationService(repositories.NewClueRepository()),
		BlockService: services.NewBlockService(
			repositories.NewBlockRepository(),
			repositories.NewBlockStateRepository(),
		),
	}
}

// GetUserFromContext retrieves the user from the context
// User will always be in the context because the middleware
func (h AdminHandler) UserFromContext(ctx context.Context) *models.User {
	return ctx.Value(contextkeys.UserKey).(*models.User)
}

func (h *AdminHandler) handleError(w http.ResponseWriter, r *http.Request, logMsg string, flashMsg string, params ...interface{}) {
	h.Logger.Error(logMsg, params...)
	err := templates.Toast(*flash.NewError(flashMsg)).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error(logMsg+" - rendering template", "error", err)
	}
}

func (h *AdminHandler) handleSuccess(w http.ResponseWriter, r *http.Request, flashMsg string) {
	err := templates.Toast(*flash.NewSuccess(flashMsg)).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering success template", "error", err)
	}
}

// redirect is a helper function to redirect the user to a new page
// It accounts for htmx requests and redirects the user to the referer
func (h AdminHandler) redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		return
	}
	http.Redirect(w, r, path, http.StatusFound)
}
