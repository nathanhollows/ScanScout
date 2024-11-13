package handlers

import (
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/templates/players"
)

func (h *PlayerHandler) Next(w http.ResponseWriter, r *http.Request) {
	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		h.redirect(w, r, "/play")
		return
	}

	locations, err := h.GameplayService.SuggestNextLocations(r.Context(), team)
	if err != nil {
		h.handleError(w, r, "Next: getting next locations", "Error getting next locations", "Could not load data", err)
		return
	}

	// data["notifications"], _ = h.NotificationService.GetNotifications(r.Context(), team.Code)
	c := templates.Next(*team, locations)
	err = templates.Layout(c, "Next stops").Render(r.Context(), w)
}
