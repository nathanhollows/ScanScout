package handlers

import (
	"net/http"
	"time"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/helpers"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// StartGame starts the game immediately
func (h *AdminHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	response := h.GameManagerService.StartGame(r.Context(), user)
	for _, msg := range response.FlashMessages {
		err := templates.GameScheduleStatus(user.CurrentInstance, msg).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("StartGame: rendering template", "error", err)
		}
		return
	}
	if response.Error != nil {
		templates.Toast(*flash.NewError("Error starting game")).Render(r.Context(), w)
		h.Logger.Error("starting game", "err", response.Error)
	}
}

// StopGame stops the game immediately
func (h *AdminHandler) StopGame(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	user := h.UserFromContext(r.Context())

	response := h.GameManagerService.StopGame(r.Context(), user)
	for _, msg := range response.FlashMessages {
		err := templates.GameScheduleStatus(user.CurrentInstance, msg).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("StopGame: rendering template", "error", err)
		}
		return
	}
	if response.Error != nil {
		templates.Toast(*flash.NewError("Error starting game")).Render(r.Context(), w)
		h.Logger.Error("stopping game", "err", response.Error)
	}
}

// ScheduleGame schedules the game to start and/or end at a specific time
func (h *AdminHandler) ScheduleGame(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	user := h.UserFromContext(r.Context())

	r.ParseForm()

	var sTime, eTime time.Time
	var err error
	if r.Form.Get("set_start") != "" {
		startDate := r.Form.Get("utc_start_date")
		startTime := r.Form.Get("utc_start_time")
		sTime, err = helpers.ParseDateTime(startDate, startTime)
		if err != nil {
			h.Logger.Error("parsing start date and time", "err", err)
			err := templates.Toast(*flash.NewError("Error parsing start date and time")).Render(r.Context(), w)
			if err != nil {
				h.Logger.Error("ScheduleGame: rendering template", "error", err)
			}
			return
		}
	}

	if r.Form.Get("set_end") != "" {
		endDate := r.Form.Get("utc_end_date")
		endTime := r.Form.Get("utc_end_time")
		eTime, err = helpers.ParseDateTime(endDate, endTime)
		if err != nil {
			h.Logger.Error("parsing end date and time", "err", err)
			err := templates.Toast(*flash.NewError("Error parsing end date and time")).Render(r.Context(), w)
			if err != nil {
				h.Logger.Error("ScheduleGame: rendering template", "error", err)
			}
			return
		}
	}

	if sTime.After(eTime) && !eTime.IsZero() {
		err := templates.Toast(*flash.NewError("Start time must be before end time")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("ScheduleGame: rendering template", "error", err)
		}
		return
	}

	response := h.GameManagerService.ScheduleGame(r.Context(), user, sTime, eTime)
	for _, msg := range response.FlashMessages {
		err := templates.GameScheduleStatus(user.CurrentInstance, msg).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("ScheduleGame: rendering template", "error", err)
		}
		return
	}
	if response.Error != nil {
		templates.Toast(*flash.NewError("Error scheduling game")).Render(r.Context(), w)
		h.Logger.Error("scheduling game", "err", response.Error)
		return
	}

	err = templates.GameScheduleStatus(user.CurrentInstance, *flash.NewSuccess("Schedule updated!")).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("ScheduleGame: rendering template", "error", err)
	}

}
