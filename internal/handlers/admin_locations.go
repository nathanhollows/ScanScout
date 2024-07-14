package handlers

import (
	"log/slog"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Locations shows admin the geosites
func AdminLocationsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Locations"

	locations, err := models.FindAllLocations(r.Context())
	if err != nil {
		flash.NewError("Error finding locations").Save(w, r)
		return
	} else {
		data["locations"] = locations
	}

	data["messages"] = flash.Get(w, r)
	render(w, data, true, "locations_index")
}

// LocationEdit shows the form to edit a location
func AdminLocationEditHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Edit Location"
	data["messages"] = flash.Get(w, r)

	// Get the location from the chi context
	code := chi.URLParam(r, "id")

	location, err := models.FindLocationByCode(r.Context(), code)
	if err != nil {
		flash.NewError("Location could not be found").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	data["location"] = location
	render(w, data, true, "locations_edit")
}

// AdminLocationEditPostHandler handles saving a location
func AdminLocationEditPostHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

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

	location, err := models.FindInstanceLocationById(r.Context(), locationCode)
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

	flash.NewSuccess("Location saved successfully").Save(w, r)
	http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
}

// LocationNew shows the form to create a new location
func AdminLocationNewHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Add a Location"
	data["messages"] = flash.Get(w, r)

	// Render the template
	render(w, data, true, "locations_new")
}

// AdminLocationNewPostHandler creates a new location
func AdminLocationNewPostHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	r.ParseForm()

	user := r.Context().Value(contextkeys.UserIDKey).(*models.User)

	name := r.FormValue("name")
	content := r.FormValue("content")
	criteriaID := r.FormValue("criteria")
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")

	err := gameManagerService.CreateLocation(r.Context(), user, name, content, criteriaID, lat, lng)
	if err != nil {
		flash.NewError("Location could not be saved").Save(w, r)
		log.Error(err, "ctx", r.Context())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Location saved").Save(w, r)
	http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
}

func adminGenerateQRHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "id")

	location, err := models.FindLocationByCode(r.Context(), code)
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
	archive, err := models.GenerateQRCodeArchive(r.Context())
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
	posters, err := models.GeneratePosters(r.Context())
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
