package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Show the form to edit the navigation settings.
func (h *AdminHandler) Navigation(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Navigation"
	data["page"] = "navigation"

	user := h.UserFromContext(r.Context())
	data["user"] = user

	data["navigation_modes"] = models.GetNavigationModes()
	data["navigation_methods"] = models.GetNavigationMethods()
	data["completion_methods"] = models.GetCompletionMethods()

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "navigation")
}

// Update the navigation settings.
func (h *AdminHandler) NavigationPost(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/navigation", http.StatusSeeOther)
		return
	}

	// Update the navigation settings
	response := h.GameManagerService.UpdateSettings(r.Context(), &user.CurrentInstance.Settings, r.Form)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		slog.Error("updating instance settings", "err", response.Error.Error())
		http.Redirect(w, r, "/admin/navigation", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/navigation", http.StatusSeeOther)
}
