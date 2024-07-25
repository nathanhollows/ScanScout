package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Show the form to edit the experience settings.
func (h *AdminHandler) Experience(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Settings"
	data["page"] = "experience"

	user := h.UserFromContext(r.Context())
	data["user"] = user

	data["navigation_modes"] = models.GetNavigationModes()
	data["navigation_methods"] = models.GetNavigationMethods()
	data["completion_methods"] = models.GetCompletionMethods()

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "experience")
}

// Update the experience settings.
func (h *AdminHandler) ExperiencePost(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Settings"
	data["page"] = "experience"

	user := h.UserFromContext(r.Context())
	data["user"] = user

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/experience", http.StatusSeeOther)
		return
	}

	// Update the experience settings
	response := h.GameManagerService.UpdateSettings(r.Context(), &user.CurrentInstance.Settings, r.Form)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		slog.Error("updating instance settings", "err", response.Error.Error())
		http.Redirect(w, r, "/admin/experience", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/experience", http.StatusSeeOther)
}
