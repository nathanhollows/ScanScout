package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/public"
)

type PublicHandler struct {
	Logger       *slog.Logger
	AuthService  services.AuthService
	EmailService services.EmailService
	UserService  services.UserService
}

func NewPublicHandler(
	logger *slog.Logger,
	authService services.AuthService,
	emailService services.EmailService,
	userService services.UserService,
) *PublicHandler {
	return &PublicHandler{
		Logger:       logger,
		AuthService:  authService,
		EmailService: emailService,
		UserService:  userService,
	}
}

func (h *PublicHandler) handleError(w http.ResponseWriter, r *http.Request, logMsg string, flashMsg string, params ...interface{}) {
	h.Logger.Error(logMsg, params...)
	err := templates.Toast(*flash.NewError(flashMsg)).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error(logMsg+" - rendering template", "error", err)
	}
}

func (h *PublicHandler) handleSuccess(w http.ResponseWriter, r *http.Request, flashMsg string) {
	err := templates.Toast(*flash.NewSuccess(flashMsg)).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering success template", "error", err)
	}
}

// redirect is a helper function to redirect the user to a new page
// It accounts for htmx requests and redirects the user to the referer
func (h *PublicHandler) redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", path)
		return
	}
	http.Redirect(w, r, path, http.StatusFound)
}
