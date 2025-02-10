package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/sessions"
	templates "github.com/nathanhollows/Rapua/internal/templates/players"
	"github.com/nathanhollows/Rapua/models"
)

// Play shows the player the first page of the game.
func (h *PlayerHandler) Play(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r, "scanscout")
	teamCode := session.Values["team"]
	var team *models.Team
	var err error

	// If the team is already playing, redirect to the next page
	if teamCode != nil {
		team, err = h.GameplayService.GetTeamByCode(r.Context(), teamCode.(string))
		if err == nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/next")
			} else {
				http.Redirect(w, r, "/next", http.StatusFound)
			}
			return
		} else {
			h.Logger.Error("Home get team from session code", "err", err, "team", teamCode)
			session.Options.MaxAge = -1
			session.Save(r, w)
		}
	}

	if team == nil {
		team = &models.Team{}
	}
	c := templates.Home(*team)
	err = templates.Layout(c, "Home", nil).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Home: rendering template", "error", err)
	}

	// Destroy the session now
	session.Options.MaxAge = -1
	session.Save(r, w)

}

// PlayPost is the handler for the play form submission.
func (h *PlayerHandler) PlayPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	teamCode := r.FormValue("team")
	teamName := r.FormValue("customTeamName")

	response := h.GameplayService.StartPlaying(r.Context(), teamCode, teamName)
	if response.Error != nil {
		err := templates.Toast(response.FlashMessages...).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("HomePost: rendering template", "error", err)
			return
		}
		return
	}

	team := response.Data["team"].(*models.Team)

	err := h.startSession(w, r, team.Code)
	if err != nil {
		h.handleError(w, r, "HomePost: starting session", "Error starting session. Please try again.", "error", err, "team", team.Code)
		return
	}

	w.Header().Set("HX-Redirect", "/next")
}
