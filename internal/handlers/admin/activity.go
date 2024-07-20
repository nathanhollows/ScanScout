package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// AdminDashboard shows the admin dashboard
func (h *AdminHandler) Activity(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Activity tracker"
	data["page"] = "activity"

	user := h.UserFromContext(r.Context())
	data["locations"] = user.CurrentInstance.Locations
	data["teams"] = user.CurrentInstance.Teams

	data["messages"] = flash.Get(w, r)
	// Render the template
	handlers.Render(w, data, handlers.AdminDir, "activity")
}
