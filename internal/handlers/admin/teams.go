package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

func (h *AdminHandler) Teams(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Teams"
	data["page"] = "teams"

	user := h.UserFromContext(r.Context())
	data["teams"] = user.CurrentInstance.Teams

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "teams_index")
}

func (h *AdminHandler) TeamsAdd(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "New Team"
	data["page"] = "teams"

	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		flash.NewError("Something went wrong behind the scenes, please try again.").Save(w, r)
		slog.Error("TeamsAdd parsing form", "error", err.Error(), "instance_id", user.CurrentInstanceID)
		http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
		return
	}

	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		flash.NewError("Something went wrong behind the scenes, please try again.").Save(w, r)
		slog.Error("TeamsAdd converting string to int", "error", err.Error(), "instance_id", user.CurrentInstanceID, "count", countStr)
		http.Redirect(w, r, "/admin/teams", http.StatusSeeOther)
		return
	}

	// Add the teams
	response := h.GameManagerService.AddTeams(r.Context(), user.CurrentInstanceID, count)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		slog.Error("TeamsAdd", "error", response.Error.Error(), "instance_id", user.CurrentInstanceID, "count", count)
	}

	http.Redirect(w, r, "/admin/teams/", http.StatusSeeOther)
}
