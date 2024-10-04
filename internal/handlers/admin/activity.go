package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Activity displays the activity tracker page
func (h *AdminHandler) Activity(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	err := h.GameManagerService.LoadTeams(r.Context(), &user.CurrentInstance.Teams)
	if err != nil {
		h.handleError(w, r, "Activity: loading teams", "Error loading teams", "Could not load data", err)
		return
	}

	c := templates.ActivityTracker(*user)
	err = templates.Layout(c, *user, "Activity", "Activity").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Activity: rendering template", "error", err)
	}
}

// ActivityTeamsOverview displays the activity tracker page
func (h *AdminHandler) ActivityTeamsOverview(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	err := h.GameManagerService.LoadTeams(r.Context(), &user.CurrentInstance.Teams)
	if err != nil {
		h.handleError(w, r, "ActivityTeamsOverview: loading teams", "Error loading teams", "Could not load data", err)
		return
	}

	err = templates.ActivityTeamsTable(user.CurrentInstance.Locations, user.CurrentInstance.Teams).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("ActivityTeamsOverview: rendering template", "error", err)
	}
}

// TeamActivity displays the activity tracker page
// It accepts HTMX requests to update the team activity
func (h *AdminHandler) TeamActivity(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	teamCode := chi.URLParam(r, "teamCode")

	gameplayService := services.NewGameplayService()
	team, err := gameplayService.GetTeamByCode(r.Context(), teamCode)
	if err != nil || team.InstanceID != user.CurrentInstanceID {
		h.handleError(w, r, "TeamActivity: getting team", "Error getting team", "Could not load data", err)
		return
	}

	team.LoadScans(r.Context())
	response := gameplayService.SuggestNextLocations(r.Context(), team, user.CurrentInstance.Settings.MaxNextLocations)
	if response.Error != nil {
		h.handleError(w, r, "TeamActivity: getting next locations", "Error getting next locations", "Could not load data", response.Error)
		return
	}

	notifications, err := h.NotificationService.GetNotifications(r.Context(), team.Code)
	if err != nil {
		h.handleError(w, r, "TeamActivity: getting notifications", "Error getting notifications", "Could not load data", err)
		return
	}

	var nextLocations models.Locations
	if response.Data["nextLocations"] == nil {
		nextLocations = models.Locations{}
	} else {
		nextLocations = response.Data["nextLocations"].(models.Locations)
	}

	err = templates.TeamActivity(user.CurrentInstance.Settings, *team, notifications, nextLocations).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("TeamActivity: rendering template", "error", err)
	}

}
