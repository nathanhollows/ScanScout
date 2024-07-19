package handlers

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

func (h *PlayerHandler) Next(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)

	session, _ := sessions.Get(r, "scanscout")
	teamCode := ""

	if r.Method == http.MethodPost {
		r.ParseForm()
		teamCode = strings.ToUpper(r.Form.Get("team"))
	} else {
		code := session.Values["team"]
		if code != nil {
			teamCode = strings.ToUpper(code.(string))
		}
	}

	if teamCode == "" {
		flash.NewError("No team code found. Please enter your team code and try again.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		session.Values["team"] = teamCode
		session.Save(r, w)
	}

	team, err := h.GameplayService.GetTeamByCode(r.Context(), teamCode)
	if err != nil {
		log.Error(err)
		flash.NewError("Team not found. Please enter your team code and try again.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if team.MustScanOut != "" {
		flash.NewInfo("You are already scanned in. You must scan out of "+team.BlockingLocation.Name+" before you can scan in to your next location.").Save(w, r)
	}

	locations, err := h.GameplayService.SuggestNextLocations(r.Context(), team, 3)
	if err != nil {
		log.Error(err)
		flash.NewError("Error suggesting next locations. Please try again later.").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data["team"] = team
	data["locations"] = locations
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, false, "next")
}
