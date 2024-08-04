package handlers

import (
	"log/slog"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

var gameManagerService *services.GameManagerService

// LocationEdit shows the form to edit a location
func AdminLocationEditHandler(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	data := TemplateData(r)
	data["title"] = "Edit Location"
	data["page"] = "locations"
	data["messages"] = flash.Get(w, r)

	// Get the location from the chi context
	code := chi.URLParam(r, "id")

	user := r.Context().Value(contextkeys.UserIDKey).(*models.User)

	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, code)
	if err != nil {
		flash.NewError("Location could not be found").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	location.LoadClues(r.Context())
	data["location"] = location
	Render(w, data, AdminDir, "locations_edit")
}

// AdminLocationEditPostHandler handles saving a location
func AdminLocationEditPostHandler(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	locationCode := chi.URLParam(r, "id")

	location, err := models.FindLocationByInstanceAndCode(
		r.Context(),
		r.Context().Value(contextkeys.UserIDKey).(*models.User).CurrentInstanceID,
		locationCode,
	)
	if err != nil {
		slog.Error("Error finding location", "err", err)
		flash.NewError("Location not found").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	newName := r.FormValue("name")
	newContent := r.FormValue("content")
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")

	err = gameManagerService.UpdateLocation(r.Context(), location, newName, newContent, lat, lng)
	if err != nil {
		log.Error(err)
		flash.NewError("Error saving location: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	err = gameManagerService.UpdateClues(r.Context(), location, r.Form["clues[]"], r.Form["clue_ids[]"])
	if err != nil {
		log.Error(err)
		flash.NewError("Could not save clues. Please try again.").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Location saved successfully").Save(w, r)
	http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
}

func adminGenerateQRHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "id")

	user := r.Context().Value(contextkeys.UserIDKey).(*models.User)

	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, code)
	if err != nil {
		flash.NewError("Location could not be found").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	err = location.Marker.GenerateQRCode()
	if err != nil {
		flash.NewError("QR code could not be generated").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}
}

func AdminLocationQRZipHandler(w http.ResponseWriter, r *http.Request) {
	archive, err := models.GenerateQRCodeArchive(r.Context(), r.Context().Value(contextkeys.UserIDKey).(*models.User).CurrentInstanceID)
	if err != nil {
		flash.NewError("QR codes could not be zipped").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	// Overwrite the file download header
	instance := r.Context().Value(contextkeys.UserIDKey).(*models.User).CurrentInstance
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" QR codes .zip")
	// Serve the file
	http.ServeFile(w, r, archive)

}

func AdminLocationPostersHandler(w http.ResponseWriter, r *http.Request) {
	posters, err := models.GeneratePosters(r.Context(), r.Context().Value(contextkeys.UserIDKey).(*models.User).CurrentInstanceID)
	if err != nil {
		log.Error(err)
		flash.NewError("Posters could not be generated").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	// Overwrite the file download header
	instance := r.Context().Value(contextkeys.UserIDKey).(*models.User).CurrentInstance
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" posters.pdf")
	// Serve the file
	http.ServeFile(w, r, posters)
}
