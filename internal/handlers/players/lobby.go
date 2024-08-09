package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Lobby is where teams wait for the game to begin
func (h *PlayerHandler) Lobby(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		flash.NewError("Error loading team.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// If the team is loaded, double check that the game is not in progress
	response := h.GameplayService.CheckGameStatus(r.Context(), team)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		slog.Error("Lobby: checking game status", "error", response.Error.Error())
		http.Redirect(w, r, "/", http.StatusFound)
	}

	status, ok := response.Data["status"].(models.GameStatus)
	if !ok {
		slog.Error("Lobby: checking game status", "error", "status not found")
		flash.NewError("Error checking game status.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	switch status {
	case models.Active:
		flash.NewSuccess("The game has started!").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	case models.Closed:
		flash.NewSuccess("The game has finished. Thanks for playing!").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		data["title"] = "Lobby"
		data["notifications"], _ = h.NotificationService.GetNotifications(r.Context(), team.Code)
		data["messages"] = flash.Get(w, r)
		team.LoadNotifications(r.Context())
		data["team"] = team
		handlers.Render(w, data, handlers.PlayerDir, "lobby")
	}
}
