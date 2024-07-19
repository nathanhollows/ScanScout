package handlers

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// PublicMyLocationsHandler shows the found locations page
func PublicMyLocationsHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData(r)
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
	Render(w, data, false, "mylocations")
}

// PublicSpecificLocationsHandler shows the page for a specific location
func PublicSpecificLocationsHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData(r)
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)
	data["code"] = locationCode

	// Get the list of locations from the session
	locations := GetLocationsFromSession(r)
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
	user := r.Context().Value(contextkeys.UserIDKey).(*models.User)
	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, locationCode)
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
	data["title"] = location.Marker.Name
	data["messages"] = flash.Get(w, r)
	Render(w, data, false, "location")
}

func GetLocationsFromSession(r *http.Request) []string {
	session, err := sessions.Get(r, "scanscout")
	if err != nil || session.Values["locations"] == nil {
		return nil
	}
	return session.Values["locations"].([]string)
}
