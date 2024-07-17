package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// AdminDashboard shows the admin dashboard
func (h *AdminHandler) Activity(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Activity tracker"
	data["breadcrumbs"] = []map[string]string{
		{"link": "/admin", "text": "Admin"},
		{"link": "/admin/dashboard", "text": "Dashboard"},
	}

	// Get the list of locations
	locations, err := h.GameManagerService.GetAllLocations(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["locations"] = locations

	// Get the list of teams and their teams
	teams, err := h.GameManagerService.GetAllTeams(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["teams"] = teams

	// Render the template
	if err := handlers.Render(w, data, true, "activity"); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
