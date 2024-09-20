package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	templates "github.com/nathanhollows/Rapua/internal/templates/players"
)

// Lobby is where teams wait for the game to begin
func (h *PlayerHandler) Lobby(w http.ResponseWriter, r *http.Request) {

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		flash.NewError("Error loading team.").Save(w, r)
		h.redirect(w, r, "/play")
		return
	}

	c := templates.Lobby(*team)
	err = templates.Layout(c, "Lobby").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering lobby", "error", err.Error())
	}
}
