package handlers

import (
	"log/slog"
	"net/http"

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
		team, err = h.GameplayService.GetTeamByCode(r.Context(), teamCode.(string))
		if err == nil {
			data["team"] = team
			http.Redirect(w, r, "/next", http.StatusFound)
			return
		} else {
			slog.Error("Home get team from session code", "err", err, "team", teamCode)
		}
	}

	// Start the game if the form is submitted
	if r.Method == http.MethodPost {
		r.ParseForm()

		response := h.GameplayService.StartPlaying(r.Context(), r.FormValue("team"), r.FormValue("customTeamName"))
		for _, message := range response.FlashMessages {
			message.Save(w, r)
		}
		if response.Error != nil {
			slog.Error("Error starting game", "err", response.Error.Error(), "team", teamCode)
		} else {
			session.Values["team"] = team.Code
			session.Save(r, w)
			data["team"] = team
			http.Redirect(w, r, "/next", http.StatusFound)
			return
		}
	}

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "home")
}
