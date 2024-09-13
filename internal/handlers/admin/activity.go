package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Activity displays the activity tracker page
func (h *AdminHandler) Activity(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.ActivityTracker(*user)
	err := templates.Layout(c, *user, "Activity", "Activity").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Activity: rendering template", "error", err)
	}
}

// ActivityTeamsOverview displays the activity tracker page
func (h *AdminHandler) ActivityTeamsOverview(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	err := templates.ActivityTeamsTable(user.CurrentInstance.Locations, user.CurrentInstance.Teams).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("ActivityTeamsOverview: rendering template", "error", err)
	}
}

// TeamActivity displays the activity tracker page
// It accepts HTMX requests to update the team activity
func (h *AdminHandler) TeamActivity(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)

	user := h.UserFromContext(r.Context())

	gameplayService := &services.GameplayService{}
	team, err := gameplayService.GetTeamByCode(r.Context(), chi.URLParam(r, "teamCode"))
	if err != nil || team.InstanceID != user.CurrentInstanceID {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}
	team.LoadScans(r.Context())
	response := gameplayService.SuggestNextLocations(r.Context(), team, user.CurrentInstance.Settings.MaxNextLocations)
	if response.Error != nil {
		http.Error(w, response.Error.Error(), http.StatusInternalServerError)
		return
	}

	notifications, err := h.NotificationService.GetNotifications(r.Context(), team.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["notifications"] = notifications
	data["settings"] = user.CurrentInstance.Settings
	data["locations"] = response.Data["nextLocations"].(models.Locations)
	data["team"] = team
	handlers.RenderHTMX(w, data, handlers.AdminDir, "team_activity")
}
