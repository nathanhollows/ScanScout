package handlers

import (
	"net/http"
)

// Dashboard shows the admin dashboard
func adminActivityHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Activity tracker"

	// Render the template
	render(w, data, true, "activity")
}
