package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"golang.org/x/exp/slog"
)

// CheckIn handles the GET request for scanning a location
func (h *PlayerHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)
	data["code"] = code

	team, ok := r.Context().Value(contextkeys.TeamKey).(*models.Team)
	if ok {
		data["team"] = team
	}
	data["team"] = team

	if team.MustScanOut != "" {
		if code == "" {
			flash.NewWarning("Please scan out at the location you scanned in at.").
				SetTitle("You are already scanned in.").Save(w, r)
			data["blocked"] = true
		} else if code == team.MustScanOut {
			message := fmt.Sprintf("Do you want to <a href=\"/o/%s\" class=\"link\">scan out</a> instead?", code)
			flash.NewDefault(message).
				SetTitle("You are already scanned in.").
				Save(w, r)
			data["blocked"] = true
		} else {
			flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).
				SetTitle("You are already scanned in.").
				Save(w, r)
			data["blocked"] = true
		}
	}

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "scan")
}

// ScanPost handles the POST request for scanning in
func (h *PlayerHandler) ScanPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	team, ok := r.Context().Value(contextkeys.TeamKey).(*models.Team)
	if !ok {
		flash.NewWarning("Team not found. Please double check the code and try again.").
			SetTitle("Team code not found").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	response := h.GameplayService.LogScan(r.Context(), team, locationCode)
	for _, msg := range response.FlashMessages {
		msg.Save(w, r)
	}
	if response.Error != nil {
		slog.Debug("Error logging scan", "err", response.Error.Error(), "team", team.Code, "location", locationCode)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	flash.NewSuccess("You have scanned in.").Save(w, r)

	session, _ := sessions.Get(r, "scanscout")
	if session.Values["locations"] == nil {
		session.Values["locations"] = []string{locationCode}
	} else {
		session.Values["locations"] = append(session.Values["locations"].([]string), locationCode)
	}
	session.Values["team"] = team.Code
	session.Save(r, w)

	http.Redirect(w, r, "/mylocations/"+locationCode, http.StatusFound)
}

func (h *PlayerHandler) ScanOut(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)

	teamCode := ""
	session, _ := sessions.Get(r, "scanscout")
	tcode := session.Values["team"]
	if tcode != nil {
		teamCode = strings.ToUpper(tcode.(string))
	}

	team, err := h.GameplayService.GetTeamByCode(r.Context(), teamCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Team code not found").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data["team"] = team

	if team.MustScanOut == "" {
		flash.NewDefault("You don't need to scan out.").
			SetTitle("You're all set!").Save(w, r)
		data["blocked"] = true
	} else if team.MustScanOut != code {
		flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).
			SetTitle("You are scanned in elsewhere.").Save(w, r)
		data["blocked"] = true
	}

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.PlayerDir, "scanout")
}

func (h *PlayerHandler) ScanOutPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	teamCode := r.FormValue("team")
	teamCode = strings.ToUpper(teamCode)

	err := h.GameplayService.LogScanOut(r.Context(), teamCode, locationCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Couldn't scan out.").Save(w, r)
		log.Error(err)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	flash.NewSuccess("You have scanned out.").Save(w, r)
	http.Redirect(w, r, "/next", http.StatusFound)
}
