package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/models"
	templates "github.com/nathanhollows/Rapua/internal/templates/players"
)

func (h *PlayerHandler) Next(w http.ResponseWriter, r *http.Request) {
	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		h.redirect(w, r, "/play")
		return
	}

	if team.MustScanOut != "" {
		h.redirect(w, r, "/checkins")
		return
	}

	response := h.GameplayService.SuggestNextLocations(r.Context(), team, 3)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("suggesting next locations", "error", response.Error.Error())
		h.redirect(w, r, "/play")
	}

	locations := response.Data["nextLocations"].(models.Locations)

	// data["notifications"], _ = h.NotificationService.GetNotifications(r.Context(), team.Code)
	c := templates.Next(*team, locations)
	err = templates.Layout(c, "Next stops").Render(r.Context(), w)
}
