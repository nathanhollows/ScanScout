package handlers

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

// publicMyLocationsHandler shows the found locations page
func publicMyLocationsHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)
	data["title"] = "My Locations"

	// Get the team code from the session
	session, _ := sessions.Get(r, "scanscout")
	teamCode := ""
	tcode := session.Values["team"]
	if tcode != nil {
		teamCode = strings.ToUpper(tcode.(string))
	}
	var team *models.Team
	var err error
	if teamCode != "" {
		team, err = models.FindTeamByCode(r.Context(), teamCode)
		if err == nil {
			data["team"] = team
		} else {
			log.Error(err)
		}
	}

	if team == nil || len(team.Scans) == 0 {
		flash.Message{
			Style:   "danger",
			Title:   "No locations found.",
			Message: "Please scan in at least one location.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data["messages"] = flash.Get(w, r)
	render(w, data, false, "mylocations")
}

// publicSpecificLocationsHandler shows the page for a specific location
func publicSpecificLocationsHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)
	data["code"] = locationCode

	// Get the list of locations from the session
	locations := getLocationsFromSession(r)
	if locations == nil {
		flash.Message{
			Style:   "danger",
			Title:   "No locations found.",
			Message: "Please scan in at least one location.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if the location is in the list
	found := false
	for _, code := range locations {
		if code == locationCode {
			found = true
			break
		}
	}
	if !found {
		flash.Message{
			Style:   "danger",
			Title:   "Location not found.",
			Message: "Please scan in at this location first.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Get the location
	location, err := models.FindLocationByCode(r.Context(), locationCode)
	if err != nil {
		flash.Message{
			Style:   "danger",
			Title:   "Location not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data["location"] = location
	data["title"] = location.Name
	data["messages"] = flash.Get(w, r)
	render(w, data, false, "location")
}

func getLocationsFromSession(r *http.Request) []string {
	session, err := sessions.Get(r, "scanscout")
	if err != nil || session.Values["locations"] == nil {
		return nil
	}
	return session.Values["locations"].([]string)
}
