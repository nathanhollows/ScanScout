package handlers

import (
	"net/http"
	"strconv"

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
		flash.Message{
			Title:   "Error",
			Message: "Locations could not be retrieved",
			Style:   flash.Error,
		}.Save(w, r)
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
		flash.Message{
			Title:   "Error",
			Message: "Location could not be found",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	data["location"] = location
	render(w, data, true, "locations_edit")
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

	var err error

	// Create a new InstanceLocation
	location := &models.Location{
		InstanceID: user.CurrentInstanceID,
		CriteriaID: r.FormValue("criteria"),
	}

	// Create the Content
	content := models.LocationContent{}
	content.Content = r.FormValue("content")
	err = content.Save(r.Context())
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Content could not be saved",
			Style:   flash.Error,
		}.Save(w, r)
		log.Error(err, "ctx", r.Context(), "content", content)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}
	location.ContentID = content.ID

	// Either a location or coordinates are required
	if !r.Form.Has("locationCode") && (!r.Form.Has("longitude") || !r.Form.Has("latitude")) {
		flash.Message{
			Title:   "Error",
			Message: "Location or coordinates are required",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	if r.Form.Has("coordsID") {
		location.CoordsID = r.FormValue("coordsID")
	}

	// Parse coordinates if location is enabled
	var lng, lat float64
	if r.FormValue("location") == "on" {
		lng, err = strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			flash.Message{
				Title:   "Error",
				Message: "Invalid coordinates",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
			return
		}
		lat, err = strconv.ParseFloat(r.FormValue("latitude"), 64)
		if err != nil {
			flash.Message{
				Title:   "Error",
				Message: "Invalid coordinates",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
			return
		}
	}

	// Create a new coords
	coords := &models.Coords{
		Name: r.FormValue("name"),
		Lat:  lat,
		Lng:  lng,
	}
	err = coords.Save(r.Context())
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Location could not be saved",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}
	location.CoordsID = coords.Code

	// Save the InstanceLocation
	err = location.Save(r.Context())
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Location could not be saved",
			Style:   flash.Error,
		}.Save(w, r)
		log.Error(err, "ctx", r.Context(), "instanceLocation", location)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	flash.Message{
		Title:   "Success",
		Message: "Location saved",
		Style:   flash.Success,
	}.Save(w, r)
	http.Redirect(w, r, "/admin/locations/"+location.CoordsID, http.StatusSeeOther)

}

func adminGenerateQRHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "id")

	location, err := models.FindLocationByCode(r.Context(), code)
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Location could not be found",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	err = location.Coords.GenerateQRCode()
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "QR code could not be generated",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}
}

func AdminLocationQRZipHandler(w http.ResponseWriter, r *http.Request) {
	archive, err := models.GenerateQRCodeArchive(r.Context())
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "QR codes could not be zipped",
			Style:   flash.Error,
		}.Save(w, r)
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
		flash.Message{
			Title:   "Error",
			Message: "Posters could not be generated",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	// Overwrite the file download header
	instance := r.Context().Value(contextkeys.UserIDKey).(*models.User).CurrentInstance
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" posters.pdf")
	// Serve the file
	http.ServeFile(w, r, posters)
}
