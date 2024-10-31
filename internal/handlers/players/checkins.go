package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/players"
	"github.com/nathanhollows/Rapua/models"
)

// CheckIn handles the GET request for scanning a location
func (h *PlayerHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)

	team, err := h.getTeamFromContext(r.Context())
	if err == nil {
		if team.MustCheckOut != "" {
			err := h.TeamService.LoadRelation(r.Context(), team, "BlockingLocation")
			if err != nil {
				h.Logger.Error("CheckIn: loading blocking location", "err", err)
				flash.NewError("Something went wrong. Please try again.").Save(w, r)
				http.Redirect(w, r, r.Header.Get("/next"), http.StatusFound)
				return
			}
		}
	}

	response := h.GameplayService.GetMarkerByCode(r.Context(), code)
	if response.Error != nil {
		h.redirect(w, r, "/404")
		return
	}

	marker, ok := response.Data["marker"].(*models.Marker)
	if !ok {
		h.redirect(w, r, "/404")
		return
	}

	c := templates.CheckIn(*marker, team.Code, team.BlockingLocation)
	err = templates.Layout(c, "Check In: "+marker.Name).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering checkin", "error", err.Error())
	}
}

// CheckInPost handles the POST request for scanning in
func (h *PlayerHandler) CheckInPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	team, err := h.getTeamFromContext(r.Context())
	if err != nil {
		team, err = h.GameplayService.GetTeamByCode(r.Context(), r.FormValue("team"))
		if err != nil {
			h.handleError(w, r, "CheckInPost: getting team by code", "Error checking in", "error", err, "team", r.FormValue("team"))
			return
		}
	}

	response := h.GameplayService.CheckIn(r.Context(), team, locationCode)
	if response.Error != nil {
		if response.Error == services.ErrAlreadyCheckedIn {
			err := templates.Toast(*flash.NewInfo("You have already checked in here.")).Render(r.Context(), w)
			if err != nil {
				h.Logger.Error("CheckInPost: rendering toast", "error", err)
			}
			return
		}
		h.handleError(w, r, "CheckInPost: checking in", "Error checking in", "error", response.Error, "team", team.Code, "location", locationCode)
		return
	}

	location, ok := response.Data["location"].(*models.Location)
	if !ok {
		h.handleError(w, r, "CheckInPost: getting location", "Error checking in", "error", fmt.Errorf("location not found"), "team", team.Code, "location", locationCode)
		return
	}

	h.redirect(w, r, "/checkins/"+location.MarkerID)
}

func (h *PlayerHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)

	team, err := h.getTeamFromContext(r.Context())
	if err == nil {
		if team.MustCheckOut != "" {
			err := h.TeamService.LoadRelation(r.Context(), team, "BlockingLocation")
			if err != nil {
				h.Logger.Error("CheckIn: loading blocking location", "err", err)
				// TODO: render error page
				h.redirect(w, r, "/404")
				return
			}
		}
	}

	response := h.GameplayService.GetMarkerByCode(r.Context(), code)
	if response.Error != nil {
		h.redirect(w, r, "/404")
		return
	}

	marker, ok := response.Data["marker"].(*models.Marker)
	if !ok {
		h.redirect(w, r, "/404")
		return
	}

	c := templates.CheckOut(*marker, team.Code, team.BlockingLocation)
	err = templates.Layout(c, "Check Out: "+marker.Name).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering checkin", "error", err.Error())
	}

}

func (h *PlayerHandler) CheckOutPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	teamCode := r.FormValue("team")
	teamCode = strings.ToUpper(teamCode)

	team, err := h.GameplayService.GetTeamByCode(r.Context(), teamCode)
	if err != nil {
		h.handleError(w, r, "CheckOutPost: getting team by code", "Error checking out. Please double check your team code.", "error", err, "team", teamCode)
		return
	}

	err = h.GameplayService.CheckOut(r.Context(), team, locationCode)
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			templates.Toast(*flash.NewError("Location not found. Please double check the code and try again.")).Render(r.Context(), w)
			return
		} else if errors.Is(err, services.ErrUnecessaryCheckOut) {
			templates.Toast(*flash.NewInfo("You are not checked in here.")).Render(r.Context(), w)
			return
		} else if errors.Is(err, services.ErrTeamNotAllowedToCheckOut) {
			templates.Toast(*flash.NewError("You are not checked in here.")).Render(r.Context(), w)
			return
		} else if errors.Is(err, services.ErrUnfinishedCheckIn) {
			templates.Toast(*flash.NewError("You must complete all activities before checking out.")).Render(r.Context(), w)
			return
		} else {
			h.handleError(w, r, "CheckOutPost: checking out", "Error checking out", "error", err, "team", team.Code, "location", locationCode)
		}
		return
	}

	h.redirect(w, r, "/next")
}

// MyCheckins shows the found locations page
func (h *PlayerHandler) MyCheckins(w http.ResponseWriter, r *http.Request) {
	team, err := h.getTeamFromContext(r.Context())
	if err != nil || team == nil {
		http.Redirect(w, r, "/play", http.StatusFound)
		return
	}

	err = h.TeamService.LoadRelations(r.Context(), team)
	if err != nil {
		flash.NewError("Error loading check ins.").Save(w, r)
		h.Logger.Error("loading check ins", "error", err.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	c := templates.MyCheckins(*team)
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

	err = h.TeamService.LoadRelations(r.Context(), team)
	if err != nil {
		flash.NewError("Error loading locations.").Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusFound)
		return
	}

	if team.MustCheckOut != "" {
		if team.BlockingLocation.MarkerID != locationCode {
			flash.NewDefault("You are currently checked into "+team.BlockingLocation.Name).Save(w, r)
		}
	}

	// Get the index of the location in the team's scans
	index := -1
	for i, scan := range team.CheckIns {
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

	blocks, blockStates, err := h.BlockService.GetBlocksWithStateByLocationIDAndTeamCode(r.Context(), team.CheckIns[index].Location.ID, team.Code)
	if err != nil {
		h.handleError(w, r, "CheckInView: getting blocks", "Error loading blocks", "error", err, "team", team.Code, "location", locationCode)
		return
	}

	c := templates.CheckInView(team.Instance.Settings, team.CheckIns[index], blocks, blockStates)
	err = templates.Layout(c, team.CheckIns[index].Location.Name).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("rendering checkin view", "error", err.Error())
	}

}
