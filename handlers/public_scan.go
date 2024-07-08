package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

// publicScanHandler shows the public scan page
func publicScanHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)
	data["code"] = code

	location, err := models.FindLocationByCode(code)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Location code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data["location"] = location

	session, _ := sessions.Get(r, "scanscout")
	teamCode := ""
	tcode := session.Values["team"]
	if tcode != nil {
		teamCode = strings.ToUpper(tcode.(string))
	}
	var team *models.Team
	if teamCode != "" {
		team, err = models.FindTeamByCode(r.Context(), teamCode)
		if err == nil {
			data["team"] = team
		} else {
			log.Error(err)
		}
	}

	// Check if the team must scan out
	if team != nil && team.MustScanOut != "" {
		if code == "" {
			flash.Message{
				Style:   "warning",
				Title:   "You are already scanned in.",
				Message: fmt.Sprint("You need to scan out at ", team.BlockingLocation.Name, "."),
			}.Save(w, r)
			data["blocked"] = true
		} else if code == team.MustScanOut {
			// Construct a message with a link to the scan out page
			message := fmt.Sprintf("Do you want to <a href=\"/o/%s\" class=\"link\">scan out</a> instead?", code)
			flash.Message{
				Style:   "",
				Title:   "You are already scanned in.",
				Message: message,
			}.Save(w, r)
			data["blocked"] = true
		} else {
			flash.Message{
				Style:   "warning",
				Title:   "You are already scanned in.",
				Message: fmt.Sprint("You need to scan out at ", team.BlockingLocation.Name, "."),
			}.Save(w, r)
			data["blocked"] = true
		}
	}

	data["messages"] = flash.Get(w, r)
	render(w, data, false, "scan")
}

// publicScanPostHandler logs the scan
func publicScanPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	// Get the location
	location, err := models.FindLocationByCode(locationCode)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Location code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if a team exists with the code
	teamCode := r.FormValue("team")
	teamCode = strings.ToUpper(teamCode)
	team, err := models.FindTeamByCode(r.Context(), teamCode)
	if err != nil || team == nil {
		flash.Message{
			Style:   "warning",
			Title:   "Team code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if the team must scan out
	if team.MustScanOut != "" {
		if locationCode != team.MustScanOut {
			flash.Message{
				Style:   "warning",
				Title:   "You are already scanned in elsewhere.",
				Message: fmt.Sprint("Please scan out at ", team.BlockingLocation.Name, " before scanning in."),
			}.Save(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			// Redirect to the scan out page
			flash.Message{
				Style:   "info",
				Title:   "You are already scanned in.",
				Message: "Do you want to scan out?",
			}.Save(w, r)
			http.Redirect(w, r, "/o/"+locationCode, http.StatusFound)
			return
		}
	}

	// Check if the team has already visited the location
	if team.HasVisited(location) {
		flash.Message{
			Style:   "warning",
			Title:   "You have already visited here.",
			Message: "Please choose another location to visit.",
		}.Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	}

	// Check if the location is one of the suggested locations
	suggested := team.SuggestNextLocations(3)
	found := false
	for _, l := range *suggested {
		if l.Code == locationCode {
			found = true
			break
		}
	}
	if !found {
		flash.Message{
			Style:   "warning",
			Title:   "Wrong location.",
			Message: "Please scan in at one of the following locations.",
		}.Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	}

	// Log the scan
	err = location.LogScan(r.Context(), teamCode)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Couldn't scan in.",
			Message: "Please check the codes and try again.",
		}.Save(w, r)
		log.Error(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if location.MustScanOut {
		team.MustScanOut = location.Code
		team.Update()
	}

	flash.Message{
		Style:   "success",
		Title:   "Success!",
		Message: "You have scanned in.",
	}.Save(w, r)

	// Append the location to the session
	session, _ := sessions.Get(r, "scanscout")
	if session.Values["locations"] == nil {
		session.Values["locations"] = []string{locationCode}
	} else {
		session.Values["locations"] = append(session.Values["locations"].([]string), locationCode)
	}
	session.Values["team"] = teamCode
	// session.Values["instance"] = location.InstanceID
	session.Save(r, w)

	http.Redirect(w, r, "/mylocations/"+locationCode, http.StatusFound)
}

func publicScanOutHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)

	location, err := models.FindLocationByCode(code)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Location code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data["location"] = location

	// Get the team code from the session
	session, _ := sessions.Get(r, "scanscout")
	teamCode := ""
	sessionCode := session.Values["team"]
	if sessionCode != nil {
		teamCode = sessionCode.(string)
	}
	teamCode = strings.ToUpper(teamCode)

	if teamCode == "" {
		flash.Message{
			Style:   "warning",
			Title:   "Team code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	team, err := models.FindTeamByCode(r.Context(), teamCode)
	if err == nil {
		data["team"] = team
	} else {
		log.Error(err)
	}

	// Check if team actually needs to scan out
	if team != nil {
		if team.MustScanOut == "" {
			flash.Message{
				Style:   "",
				Title:   "You're all set!",
				Message: "You don't need to scan out.",
			}.Save(w, r)
			data["blocked"] = true
		} else if team.MustScanOut != code {
			flash.Message{
				Style:   "warning",
				Title:   "You are scanned in elsewhere.",
				Message: fmt.Sprint("You need to scan out at ", team.BlockingLocation.Name, "."),
			}.Save(w, r)
			data["blocked"] = true
		}
	}

	data["messages"] = flash.Get(w, r)
	render(w, data, false, "scanout")
}

func publicScanOutPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	// Get the location
	location, err := models.FindLocationByCode(locationCode)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Location code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Get the team code from the form
	r.ParseForm()
	teamCode := r.FormValue("team")
	teamCode = strings.ToUpper(teamCode)

	team, err := models.FindTeamByCode(r.Context(), teamCode)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Team code not found.",
			Message: "Please double check the code and try again.",
		}.Save(w, r)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Check if the team must scan out
	if team.MustScanOut == "" {
		flash.Message{
			Style:   "",
			Title:   "You're all set!",
			Message: "You don't need to scan out.",
		}.Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	} else if team.MustScanOut != locationCode {
		flash.Message{
			Style:   "warning",
			Title:   "You are scanned in elsewhere.",
			Message: fmt.Sprint("You need to scan out at ", team.BlockingLocation.Name, "."),
		}.Save(w, r)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Log the scan
	err = location.LogScanOut(teamCode)
	if err != nil {
		flash.Message{
			Style:   "warning",
			Title:   "Couldn't scan out.",
			Message: "Please check the codes and try again.",
		}.Save(w, r)
		log.Error(err)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Clear the must scan out field
	team.MustScanOut = ""
	team.Update()

	flash.Message{
		Style:   "success",
		Title:   "Success!",
		Message: "You have scanned out.",
	}.Save(w, r)
	http.Redirect(w, r, "/next", http.StatusFound)

}
