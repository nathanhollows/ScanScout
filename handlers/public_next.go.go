package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

// publicNextHandler shows the team the next location(s) to scan in
func publicNextHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)

	session, _ := sessions.Get(r, "scanscout")
	teamCode := session.Values["team"]

	if teamCode == nil {
		r.ParseForm()
		teamCode = r.Form.Get("team")
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
	}

	var team *models.Team
	var err error
	if teamCode != nil {
		team, err = models.FindTeamByCode(teamCode.(string))
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

	data["locations"] = team.SuggestNextLocations(3)

	data["messages"] = flash.Get(w, r)
	render(w, data, false, "next")
}
