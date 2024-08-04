package handlers

import (
	"log/slog"
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
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if team.MustScanOut != "" {
		flash.NewInfo("You are already scanned in. You must scan out of "+team.BlockingLocation.Name+" before you can scan in to your next location.").Save(w, r)
	}

	response := h.GameplayService.SuggestNextLocations(r.Context(), team, 3)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		slog.Error("suggesting next locations", "error", response.Error.Error())
		http.Redirect(w, r, "/", http.StatusFound)
	}

	locations := response.Data["nextLocations"].(models.Locations)

	data["team"] = team
	data["locations"] = locations
	data["title"] = "Next Stop"
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "next")
}
