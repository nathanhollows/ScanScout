package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/ScanScout/models"
)

// Dashboard shows the admin dashboard
func adminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Activity tracker"

	// Get the list of locations
	locations, err := models.FindAllInstanceLocations(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["locations"] = locations

	// Get the list of teams and their activity
	activity, err := models.TeamActivityOverview(r.Context())
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["activity"] = activity

	// Render the template
	render(w, data, true, "dashboard")
}
