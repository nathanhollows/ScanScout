package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// PlayerHandler handles the player routes
func (h *PlayerHandler) Scan(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)
	data["code"] = code

	location, err := models.FindLocationByCode(r.Context(), code)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Location code note found").Save(w, r)
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
			flash.NewWarning("Please scan out at the location you scanned in at.").
				SetTitle("You are already scanned in.").Save(w, r)
			data["blocked"] = true
		} else if code == team.MustScanOut {
			// Construct a message with a link to the scan out page
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
	handlers.Render(w, data, false, "scan")
}

// ScanPost handles the POST request for scanning in
func (h *PlayerHandler) ScanPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	// Get the location
	location, err := models.FindLocationByCode(r.Context(), locationCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Location code not found").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if a team exists with the code
	teamCode := r.FormValue("team")
	teamCode = strings.ToUpper(teamCode)
	team, err := models.FindTeamByCode(r.Context(), teamCode)
	if err != nil || team == nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Team code not found").Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if the team must scan out
	if team.MustScanOut != "" {
		if locationCode != team.MustScanOut {
			flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).
				SetTitle("You are already scanned in.").
				Save(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		} else {
			// Redirect to the scan out page
			flash.NewInfo("Do you want to scan out instead?").
				SetTitle("You are already scanned in.").
				Save(w, r)
			http.Redirect(w, r, "/o/"+locationCode, http.StatusFound)
			return
		}
	}

	// Check if the team has already visited the location
	if team.HasVisited(&location.Marker) {
		flash.NewWarning("Please choose another location to visit").
			SetTitle("You have already visited here.").
			Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	}

	// Check if the location is one of the suggested locations
	suggested := team.SuggestNextLocations(r.Context(), 3)
	found := false
	for _, l := range *suggested {
		if l.Code == locationCode {
			found = true
			break
		}
	}
	if !found {
		flash.NewWarning("Please choose another location to visit").
			SetTitle("Nice try.").
			Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	}

	// Log the scan
	err = location.Marker.LogScan(r.Context(), teamCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Couldn't scan in.").
			Save(w, r)
		log.Error(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// if location.MustScanOut {
	// 	team.MustScanOut = location.Code
	// 	team.Update(r.Context())
	// }

	flash.NewSuccess("You have scanned in.").Save(w, r)

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

func (h *PlayerHandler) ScanOut(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	code := chi.URLParam(r, "code")
	code = strings.ToUpper(code)

	location, err := models.FindLocationByCode(r.Context(), code)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Location code not found").Save(w, r)
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
		flash.NewWarning("Please double check the code and try again").
			SetTitle("Team code not found").Save(w, r)
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
			flash.NewDefault("You don't need to scan out.").
				SetTitle("You're all set!").Save(w, r)
			data["blocked"] = true
		} else if team.MustScanOut != code {
			flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).
				SetTitle("You are scanned in elsewhere.").Save(w, r)
			data["blocked"] = true
		}
	}

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, false, "scanout")
}

func (h *PlayerHandler) ScanOutPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	locationCode := chi.URLParam(r, "code")
	locationCode = strings.ToUpper(locationCode)

	// Get the location
	location, err := models.FindLocationByCode(r.Context(), locationCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Location code not found").Save(w, r)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Get the team code from the form
	r.ParseForm()
	teamCode := r.FormValue("team")
	teamCode = strings.ToUpper(teamCode)

	team, err := models.FindTeamByCode(r.Context(), teamCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Team code not found").Save(w, r)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Check if the team must scan out
	if team.MustScanOut == "" {
		flash.NewWarning("You don't need to scan out.").
			SetTitle("You're all set!").Save(w, r)
		http.Redirect(w, r, "/next", http.StatusFound)
		return
	} else if team.MustScanOut != locationCode {
		flash.NewWarning(fmt.Sprintf("You need to scan out at %s.", team.BlockingLocation.Name)).
			SetTitle("You are scanned in elsewhere.").Save(w, r)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Log the scan
	err = location.Marker.LogScanOut(r.Context(), teamCode)
	if err != nil {
		flash.NewWarning("Please double check the code and try again.").
			SetTitle("Couldn't scan out.").Save(w, r)
		log.Error(err)
		http.Redirect(w, r, "/mylocations", http.StatusFound)
		return
	}

	// Clear the must scan out field
	team.MustScanOut = ""
	team.Update(r.Context())

	flash.NewSuccess("You have scanned out.").Save(w, r)
	http.Redirect(w, r, "/next", http.StatusFound)

}
