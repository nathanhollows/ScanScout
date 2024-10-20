package handlers

import (
	"net/http"

	admin "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Show the form to edit the navigation settings.
func (h *AdminHandler) Experience(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := admin.Experience(user.CurrentInstance.Settings)
	err := admin.Layout(c, *user, "Experience", "Experience").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering navigation page", "error", err.Error())
	}

}

// Update the navigation settings.
func (h *AdminHandler) ExperiencePost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.handleError(w, r, "Error parsing form", "Error parsing form", "error", err)
		return
	}

	// Update the navigation settings
	response := h.GameManagerService.UpdateSettings(r.Context(), &user.CurrentInstance.Settings, r.Form)
	if response.Error != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		h.handleError(w, r, "updating instance settings", "Error updating instance settings", "error", response.Error)
		return
	}

	err := admin.Toast(response.FlashMessages...).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering toast", "error", err.Error())
	}
}
