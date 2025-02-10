package handlers

import (
	"errors"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/services"
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
		if errors.Is(err, services.ErrAllLocationsVisited) && team.MustCheckOut == "" {
			h.redirect(w, r, "/finish")
			return
		}
		h.handleError(w, r, "Next: getting next locations", "Error getting next locations", "Could not load data", err)
		return
	}

	// data["notifications"], _ = h.NotificationService.GetNotifications(r.Context(), team.Code)
	c := templates.Next(*team, locations)
	err = templates.Layout(c, "Next stops", team.Messages).Render(r.Context(), w)
	if err != nil {
		h.handleError(w, r, "Next: rendering template", "Error rendering template", "Could not render template", err)
	}
}
