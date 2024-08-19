package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/handlers"
)

// StartGame starts the game immediately
func (h *AdminHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	user := h.UserFromContext(r.Context())

	response := h.GameManagerService.StartGame(r.Context(), user)
	for _, msg := range response.FlashMessages {
		msg.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("starting game", "err", response.Error)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// StopGame stops the game immediately
func (h *AdminHandler) StopGame(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	user := h.UserFromContext(r.Context())

	response := h.GameManagerService.StopGame(r.Context(), user)
	for _, msg := range response.FlashMessages {
		msg.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("stopping game", "err", response.Error)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// ScheduleGame schedules the game to start and/or end at a specific time
func (h *AdminHandler) ScheduleGame(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	user := h.UserFromContext(r.Context())

	r.ParseForm()

	response := h.GameManagerService.ScheduleGame(r.Context(), user, r.Form)
	for _, msg := range response.FlashMessages {
		msg.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("scheduling game", "err", response.Error)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
