package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

// publicHomeHandler shows the public home page
func publicHomeHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)

	session, _ := sessions.Get(r, "scanscout")
	teamCode := session.Values["team"]
	var team *models.Team
	var err error
	if teamCode != nil {
		team, err = models.FindTeamByCode(teamCode.(string))
		if err == nil {
			data["team"] = team
		} else {
			log.Error(err)
		}
	}

	data["messages"] = flash.Get(w, r)
	render(w, data, false, "home")
}
