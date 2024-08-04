package handlers

import (
	"log/slog"
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

	if team == nil {
		flash.NewError("You haven't started a game yet. Please enter your team code to start.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = team.LoadScans(r.Context())
	if err != nil {
		flash.NewError("Error loading check ins.").Save(w, r)
		slog.Error("loading check ins", "error", err.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	err = team.LoadBlockingLocation(r.Context())
	if err != nil {
		flash.NewError("Error loading blocking location.").Save(w, r)
		slog.Error("loading blocking location", "error", err.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	if team.MustScanOut != "" {
		flash.NewWarning("You need to scan out at "+team.BlockingLocation.Name).Save(w, r)
	}

	if len(team.Scans) == 0 {
		flash.Message{
			Style:   flash.Default,
			Message: "You haven't checked in anywhere yet.",
		}.Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
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

	err = team.LoadBlockingLocation(r.Context())
	if err != nil {
		flash.NewError("Error loading blocking location.").Save(w, r)
		slog.Error("loading blocking location", "error", err.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	if team.MustScanOut != "" {
		if team.BlockingLocation.MarkerID != locationCode {
			flash.NewDefault("You are currently checked into "+team.BlockingLocation.Name).Save(w, r)
		}
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
		if scan.Location.MarkerID == locationCode {
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
