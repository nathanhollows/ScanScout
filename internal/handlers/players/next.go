package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

func (h *PlayerHandler) Next(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)

	team := h.getTeamFromContext(r.Context())

	if team.MustScanOut != "" {
		flash.NewInfo("You are already scanned in. You must scan out of "+team.BlockingLocation.Name+" before you can scan in to your next location.").Save(w, r)
	}

	locations, err := h.GameplayService.SuggestNextLocations(r.Context(), team, 3)
	if err != nil {
		log.Error(err)
		flash.NewError("Error suggesting next locations. Please try again later.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data["team"] = team
	data["locations"] = locations
	data["title"] = "Next Locations"
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "next")
}
