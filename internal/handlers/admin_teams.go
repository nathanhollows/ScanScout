package handlers

import (
	"net/http"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Teams shows admin the teams
func AdminTeamsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Teams"

	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.NewError("User not authenticated").Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	teams, err := gameManagerService.GetAllTeams(r.Context(), user)
	if err != nil {
		flash.NewError("Error finding teams").Save(w, r)
	} else {
		data["teams"] = teams
	}

	data["messages"] = flash.Get(w, r)

	// Render the template
	render(w, data, true, "teams_index")
}

// AddTeams creates new teams equal to the number of teams in the request
func AdminTeamsAddHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	user, ok := templateData(r)["user"].(*models.User)
	if !ok || user == nil {
		flash.NewError("User not authenticated").Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		flash.NewError("Invalid number of teams").Save(w, r)
		http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
		return
	}

	// Add the teams
	err = gameManagerService.AddTeams(r.Context(), user, count)
	if err != nil {
		flash.NewError("Error adding teams: "+err.Error()).Save(w, r)
	} else {
		flash.NewSuccess("Teams added").Save(w, r)
	}

	// Redirect to the teams page
	http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
}
