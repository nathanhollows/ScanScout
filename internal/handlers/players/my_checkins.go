package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	templates "github.com/nathanhollows/Rapua/internal/templates/players"
)

// MyCheckins shows the found locations page
func (h *PlayerHandler) MyCheckins(w http.ResponseWriter, r *http.Request) {
	team, err := h.getTeamFromContext(r.Context())
	if err != nil || team == nil {
		http.Redirect(w, r, "/play", http.StatusFound)
		return
	}

	err = team.LoadScans(r.Context())
	if err != nil {
		flash.NewError("Error loading check ins.").Save(w, r)
		h.Logger.Error("loading check ins", "error", err.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	err = team.LoadBlockingLocation(r.Context())
	if err != nil {
		// We don't want to stop the user from seeing their check-ins if the blocking location can't be loaded
		h.Logger.Error("loading blocking location", "error", err.Error())
	}

	if len(team.Scans) == 0 {
		flash.Message{
			Style:   flash.Default,
			Message: "You haven't checked in anywhere yet.",
		}.Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	}

	// TODO: Handle notifications
	// notifications, _ := h.NotificationService.GetNotifications(r.Context(), team.Code)

	c := templates.Checkins(*team)
	err = templates.Layout(c, "My Check-ins").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering checkins", "error", err.Error())
	}
}

// CheckInView shows the page for a specific location
func (h *PlayerHandler) CheckInView(w http.ResponseWriter, r *http.Request) {
	locationCode := chi.URLParam(r, "id")

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		flash.NewError("Error loading team.").Save(w, r)
		http.Redirect(w, r, "/play", http.StatusFound)
		return
	}

	err = team.LoadBlockingLocation(r.Context())
	if err != nil {
		flash.NewError("Error loading blocking location.").Save(w, r)
		h.Logger.Error("loading blocking location", "error", err.Error())
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

	c := templates.CheckInView(team.Scans[index])
	err = templates.Layout(c, team.Scans[index].Location.Name).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering checkin view", "error", err.Error())
	}

}
