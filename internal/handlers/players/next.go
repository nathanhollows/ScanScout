package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
)

func (h *PlayerHandler) Next(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		flash.NewError("Error loading team.").Save(w, r)
		http.Redirect(w, r, "/play", http.StatusFound)
		return
	}

	if team.MustScanOut != "" {
		http.Redirect(w, r, "/checkins", http.StatusFound)
		return
	}

	response := h.GameplayService.SuggestNextLocations(r.Context(), team, 3)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("suggesting next locations", "error", response.Error.Error())
		http.Redirect(w, r, "/play", http.StatusFound)
	}

	if response.Data["blockingLocation"] != nil {
		data["blocked"] = true
	} else {
		locations := response.Data["nextLocations"].(models.Locations)

		data["team"] = team
		data["locations"] = locations
	}

	data["title"] = "Next Stop"
	data["messages"] = flash.Get(w, r)
	data["notifications"], _ = h.NotificationService.GetNotifications(r.Context(), team.Code)
	handlers.Render(w, data, handlers.PlayerDir, "next")
}
