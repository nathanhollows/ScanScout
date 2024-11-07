package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/repositories"
	internalServices "github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/services"
)

type AdminHandler struct {
	Logger              *slog.Logger
	GameManagerService  internalServices.GameManagerService
	NotificationService internalServices.NotificationService
	UserServices        internalServices.UserServices
	AssetGenerator      internalServices.AssetGenerator
	LocationService     internalServices.LocationService
	BlockService        internalServices.BlockService
	TeamService         internalServices.TeamService
	ClueService         internalServices.ClueService
	PlanService         services.PlanService
}

func NewAdminHandler(logger *slog.Logger, gs internalServices.GameManagerService, ns internalServices.NotificationService, us internalServices.UserServices) *AdminHandler {
	return &AdminHandler{
		Logger:              logger,
		GameManagerService:  gs,
		NotificationService: ns,
		UserServices:        us,
		AssetGenerator:      internalServices.NewAssetGenerator(),
		LocationService: internalServices.NewLocationService(
			repositories.NewClueRepository(),
		),
		BlockService: internalServices.NewBlockService(
			repositories.NewBlockRepository(),
			repositories.NewBlockStateRepository(),
		),
		TeamService: internalServices.NewTeamService(
			repositories.NewTeamRepository(),
		),
		ClueService: internalServices.NewClueService(
			repositories.NewClueRepository(),
			repositories.NewLocationRepository(),
		),
		PlanService: internalServices.NewPlanService(),
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
