package handlers

import (
	"net/http"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

func (h *AdminHandler) Teams(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Teams"

	user := h.UserFromContext(r.Context())
	data["teams"] = user.CurrentInstance.Teams

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "teams_index")
}

func (h *AdminHandler) TeamsAdd(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "New Team"

	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
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
	err = h.GameManagerService.AddTeams(r.Context(), user.CurrentInstanceID, count)
	if err != nil {
		flash.NewError("Error adding teams: "+err.Error()).Save(w, r)
	} else {
		flash.NewSuccess("Teams added").Save(w, r)
	}

	flash.NewSuccess("Team created successfully").Save(w, r)
	http.Redirect(w, r, "/admin/teams/", http.StatusSeeOther)
}
