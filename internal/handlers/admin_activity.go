package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
)

// AdminDashboard shows the admin dashboard
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Activity tracker"
	data["breadcrumbs"] = []map[string]string{
		{"link": "/admin", "text": "Admin"},
		{"link": "/admin/dashboard", "text": "Dashboard"},
	}

	// Get the list of locations
	locations, err := gameManagerService.GetAllLocations(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["locations"] = locations

	// Get the list of teams and their activity
	activity, err := gameManagerService.GetTeamActivityOverview(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["activity"] = activity

	// Render the template
	if err := render(w, data, true, "dashboard"); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
