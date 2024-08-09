package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

// Activity displays the activity tracker page
func (h *AdminHandler) Activity(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Activity tracker"
	data["page"] = "activity"

	user := h.UserFromContext(r.Context())
	data["locations"] = user.CurrentInstance.Locations
	for i := range user.CurrentInstance.Teams {
		if !user.CurrentInstance.Teams[i].HasStarted {
			continue
		}
		user.CurrentInstance.Teams[i].LoadScans(r.Context())
	}
	data["teams"] = user.CurrentInstance.Teams

	data["messages"] = flash.Get(w, r)
	// Render the template
	handlers.Render(w, data, handlers.AdminDir, "activity")
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
