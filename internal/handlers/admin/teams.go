package handlers

import (
	"net/http"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/models"
	admin "github.com/nathanhollows/Rapua/internal/templates/admin"
)

func (h *AdminHandler) Teams(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := admin.Teams(user.CurrentInstance.Teams)
	err := admin.Layout(c, *user, "Teams", "Teams").Render(r.Context(), w)

	if err != nil {
		h.Logger.Error("rendering teams page", "error", err.Error())
	}
}

func (h *AdminHandler) TeamsAdd(w http.ResponseWriter, r *http.Request) {

	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "TeamsAdd parsing form", "Error adding teams", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		h.handleError(w, r, "TeamsAdd parsing count", "Error adding teams", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	// Add the teams
	response := h.GameManagerService.AddTeams(r.Context(), user.CurrentInstanceID, count)
	if response.Error != nil {
		h.handleError(w, r, "TeamsAdd adding teams", "Error adding teams", "error", response.Error, "instance_id", user.CurrentInstanceID, "count", count)
		return
	}

	teams := response.Data["teams"].(models.Teams)
	err = admin.TeamsList(teams).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("TeamsAdd rendering teams list", "error", err.Error(), "instance_id", user.CurrentInstanceID)
	}
}
