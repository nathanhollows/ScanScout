package handlers

import (
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/flash"
	"github.com/nathanhollows/Rapua/models"
)

// Locations shows admin the geosites
func adminLocationsHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Locations"

	locations, err := models.FindAllInstanceLocations(r.Context())
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
func adminLocationEditHandler(w http.ResponseWriter, r *http.Request) {
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
func adminLocationNewHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Add a Location"
	data["messages"] = flash.Get(w, r)

	// Render the template
	render(w, data, true, "locations_new")
}

// saveLocation saves a new location
func adminLocationSaveHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	r.ParseForm()

	var lng, lat float64
	if r.FormValue("location") == "on" {
		var errLng, errLat error
		lng, errLng = strconv.ParseFloat(r.FormValue("longitude"), 32)
		lat, errLat = strconv.ParseFloat(r.FormValue("latitude"), 32)

		if errLng != nil || errLat != nil {
			flash.Message{
				Title:   "Error",
				Message: "Invalid coordinates",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/admin/locations/new", http.StatusSeeOther)
			return
		}
	}

	// Generate a code if one is not provided
	location := &models.Location{}
	var err error
	if r.Form.Has("code") {
		location, err = models.FindLocationByCode(r.Context(), r.FormValue("code"))
		if err != nil {
			flash.Message{
				Title:   "Error",
				Message: "Location could not be found",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/admin/locations/new", http.StatusSeeOther)
			return
		}
	}

	// Update the location
	location.Name = r.FormValue("name")
	location.Content = r.FormValue("content")
	location.Lat = lat
	location.Lng = lng

	// Save the location
	err = location.Save(r.Context())
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Location could not be saved",
			Style:   flash.Error,
		}.Save(w, r)

		http.Redirect(w, r, "/admin/locations/new", http.StatusSeeOther)
		return
	}

	flash.Message{
		Title:   "Success",
		Message: "Location saved",
		Style:   flash.Success,
	}.Save(w, r)

	http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
}

func adminLocationQRHandler(w http.ResponseWriter, r *http.Request) {
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

	err = location.GenerateQRCode()
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

func adminLocationQRZipHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "id")

	instance, err := models.FindInstanceByID(r.Context(), code)
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Instance could not be found",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	archive, err := instance.ZipQRCodes(r.Context())
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
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" QR codes .zip")
	// Serve the file
	http.ServeFile(w, r, archive)

}

func adminLocationPostersHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "id")

	instance, err := models.FindInstanceByID(r.Context(), code)
	if err != nil {
		log.Error(err)
		flash.Message{
			Title:   "Error",
			Message: "Instance could not be found",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	posters, err := instance.GeneratePosters(r.Context())
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
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" posters.pdf")
	// Serve the file
	http.ServeFile(w, r, posters)
}
