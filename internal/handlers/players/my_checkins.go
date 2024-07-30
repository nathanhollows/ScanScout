package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// CheckInList shows the found locations page
func (h *PlayerHandler) CheckInList(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	data["title"] = "My Check-ins"

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		flash.NewError("Error loading team.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if team == nil || len(team.Scans) == 0 {
		flash.Message{
			Style:   "danger",
			Title:   "No locations found.",
			Message: "You haven't checked in anywhere yet.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = team.LoadScans(r.Context())
	if err != nil {
		flash.NewError("Error loading locations.").Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	data["team"] = team
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "checkins")
}

// CheckInView shows the page for a specific location
func (h *PlayerHandler) CheckInView(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	locationCode := chi.URLParam(r, "id")

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		flash.NewError("Error loading team.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = team.LoadScans(r.Context())
	if err != nil {
		flash.NewError("Error loading locations.").Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	// Get the index of the location in the team's scans
	index := -1
	for i, scan := range team.Scans {
		if scan.LocationID == locationCode {
			index = i
			break
		}
	}

	if index == -1 {
		flash.NewWarning("Please double check the code and try again.").Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	data["title"] = team.Scans[index].Location.Name
	data["scan"] = team.Scans[index]
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "checkin_view")
}
