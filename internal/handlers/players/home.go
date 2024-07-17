package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// Home shows the public home page
func (h *PlayerHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)

	session, _ := sessions.Get(r, "scanscout")
	teamCode := session.Values["team"]
	var team *models.Team
	var err error

	// If the team is already playing, redirect to the next page
	if teamCode != nil {
		team, err = h.GameplayService.GetTeamStatus(r.Context(), strings.ToUpper(teamCode.(string)))
		if err == nil {
			data["team"] = team
			http.Redirect(w, r, "/next", http.StatusFound)
			return
		} else {
			slog.Error("Error getting team status", "err", err, "team", teamCode)
		}
	}

	// Start the game if the form is submitted
	if r.Method == http.MethodPost {
		r.ParseForm()
		teamCode := strings.ToUpper(r.FormValue("team"))
		customTeamName := r.FormValue("customTeamName")

		team, err := h.GameplayService.StartPlaying(r.Context(), teamCode, customTeamName)
		if err != nil {
			slog.Error("Error starting game", "err", err, "team", teamCode)
			flash.NewError(err.Error()).Save(w, r)
		} else {
			session.Values["team"] = team.Code
			session.Save(r, w)
			data["team"] = team
			flash.NewSuccess("Game started!").Save(w, r)
		}
	}

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, false, "home")
}
