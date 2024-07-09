package handlers

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/flash"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/sessions"
)

// publicNextHandler shows the team the next location(s) to scan in
func publicNextHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)

	// Get the team code
	teamCode := ""
	session, _ := sessions.Get(r, "scanscout")
	if r.Method == "POST" {
		r.ParseForm()
		teamCode = r.Form.Get("team")
	} else {
		code := session.Values["team"]
		if code != nil {
			teamCode = code.(string)
		}
	}
	teamCode = strings.ToUpper(teamCode)

	// If no team code is found, redirect to the home page
	if teamCode == "" {
		flash.Message{
			Style:   "warning",
			Title:   "No team code found.",
			Message: "Please enter your team code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		session.Values["team"] = teamCode
		session.Save(r, w)
	}

	// Get the team
	var team *models.Team
	var err error
	if teamCode != "" {
		team, err = models.FindTeamByCode(r.Context(), teamCode)
		if err == nil {
			data["team"] = team
		} else {
			log.Error(err)
			flash.Message{
				Style:   "danger",
				Title:   "Team not found.",
				Message: "Please enter your team code and try again.",
			}.Save(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	// If team is currently blocked, then show an alert
	if team.MustScanOut != "" {
		flash.Message{
			Style:   "info",
			Title:   "You are already scanned in.",
			Message: "You must scan out of " + team.BlockingLocation.Name + " before you can scan in to your next location.",
		}.Save(w, r)
	}

	data["locations"] = team.SuggestNextLocations(r.Context(), 3)
	data["messages"] = flash.Get(w, r)
	render(w, data, false, "next")
}
