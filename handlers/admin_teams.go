package handlers

import (
	"net/http"
	"strconv"

	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
)

// Teams shows admin the teams
func adminTeamsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Teams"

	data["messages"] = flash.Get(w, r)

	teams, err := models.FindAllTeams(r.Context())
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error finding teams",
			Style:   flash.Error,
		}.Save(w, r)
	} else {
		data["teams"] = teams
	}

	// Render the template
	render(w, data, true, "teams_index")
}

// AddTeams creates new teams equal to the number of teams in the request
func adminTeamsAddHandler(w http.ResponseWriter, r *http.Request) {
	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		flash.Message{
			Title:   "Error",
			Message: "Invalid number of teams",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
		return
	}

	// Add the teams
	err = models.AddTeams(r.Context(), count)
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error adding teams",
			Style:   flash.Error,
		}.Save(w, r)
	} else {
		flash.Message{
			Title:   "Success",
			Message: "Teams added",
			Style:   flash.Success,
		}.Save(w, r)
	}

	// Redirect to the teams page
	http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
}
