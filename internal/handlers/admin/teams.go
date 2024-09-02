package handlers

import (
	"net/http"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/flash"
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
		h.Logger.Error("TeamsAdd parsing form", "error", err.Error(), "instance_id", user.CurrentInstanceID)
		message := flash.NewError("Could not add teams, please try again.")
		err := admin.Toast(*message).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("TeamsAdd rendering toast", "error", err.Error())
		}
		return
	}

	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		h.Logger.Error("TeamsAdd parsing count", "error", err.Error(), "instance_id", user.CurrentInstanceID)
		err := admin.Toast(*flash.NewError("Could not add teams, please try again.")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("TeamsAdd rendering toast", "error", err.Error(), "instance_id", user.CurrentInstanceID)
		}
		return
	}

	// Add the teams
	response := h.GameManagerService.AddTeams(r.Context(), user.CurrentInstanceID, count)
	if response.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logger.Error("TeamsAdd", "error", response.Error.Error(), "instance_id", user.CurrentInstanceID, "count", count)
		err := admin.Toast(response.FlashMessages...).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("TeamsAdd rendering toast", "error", err.Error())
		}
		return
	}

	teams := response.Data["teams"].(models.Teams)
	err = admin.TeamsList(teams).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("TeamsAdd rendering teams list", "error", err.Error(), "instance_id", user.CurrentInstanceID)
	}
}
