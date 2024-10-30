package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Locations shows admin the locations
func (h *AdminHandler) Locations(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	for i, location := range user.CurrentInstance.Locations {
		err := h.LocationService.LoadRelations(r.Context(), &location)
		if err != nil {
			h.handleError(w, r, "Locations: loading relations", "Error loading relations", "error", err, "instance_id", user.CurrentInstanceID)
			return
		}
		user.CurrentInstance.Locations[i] = location
	}

	c := templates.LocationsIndex(user.CurrentInstance.Settings, user.CurrentInstance.Locations)
	err := templates.Layout(c, *user, "Locations", "Locations").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Locations: rendering template", "error", err)
	}

}

// LocationNew shows the form to create a new location
func (h *AdminHandler) LocationNew(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.AddLocation(user.CurrentInstance.Settings)
	err := templates.Layout(c, *user, "Locations", "New Location").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("LocationNew: rendering template", "error", err)
	}
}

// LocationNewPost handles creating a new location
func (h *AdminHandler) LocationNewPost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.handleError(w, r, "LocationNewPost: parsing form", "Error parsing form", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	data := make(map[string]string)
	for key, value := range r.Form {
		data[key] = value[0]
	}

	location, err := h.GameManagerService.CreateLocation(r.Context(), user, data)
	if err != nil {
		h.handleError(w, r, "LocationNewPost: creating location", "Error creating location", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	h.redirect(w, r, "/admin/locations/"+location.MarkerID)
}

// ReorderLocations handles reordering locations
// Returns a 200 status code if successful
// Otherwise, returns a 500 status code
func (h *AdminHandler) ReorderLocations(w http.ResponseWriter, r *http.Request) {
	// Check HTMX headers
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := h.UserFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.handleError(w, r, "ReorderLocations: parsing form", "Error parsing form", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	locations := r.Form["location"]
	response := h.GameManagerService.ReorderLocations(r.Context(), user, locations)
	// Discard the flash messages since this is invoked via HTMX
	if response.Error != nil {
		h.handleError(w, r, "ReorderLocations: reordering locations", "Error reordering locations", "error", response.Error, "instance_id", user.CurrentInstanceID)
		return
	}

	h.handleSuccess(w, r, "Order updated")
}

// LocationEdit shows the form to edit a location
func (h *AdminHandler) LocationEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get the location from the chi context
	code := chi.URLParam(r, "id")
	user := h.UserFromContext(r.Context())

	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, code)
	if err != nil {
		h.Logger.Error("LocationEdit: finding location", "error", err, "instance_id", user.CurrentInstanceID, "location_code", code)
		h.redirect(w, r, "/admin/locations")
		return
	}

	blocks, err := h.BlockService.GetByLocationID(r.Context(), location.ID)
	if err != nil {
		h.Logger.Error("LocationEdit: getting blocks", "error", err, "instance_id", user.CurrentInstanceID, "location_id", location.ID)
		h.redirect(w, r, "/admin/locations")
		return
	}

	err = h.LocationService.LoadCluesForLocation(r.Context(), location)
	if err != nil {
		h.handleError(w, r, "LocationEdit: loading clues", "Error loading clues", "error", err, "instance_id", user.CurrentInstanceID, "location_id", location.ID)
		return
	}

	c := templates.EditLocation(*location, user.CurrentInstance.Settings, blocks)
	err = templates.Layout(c, *user, "Locations", "Edit Location").Render(r.Context(), w)
}

// LocationEditPost handles updating a location
func (h *AdminHandler) LocationEditPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "LocationEditPost: parsing form", "Error parsing form", "error", err)
		return
	}

	user := h.UserFromContext(r.Context())
	locationCode := chi.URLParam(r, "id")

	var points int
	var err error
	if r.FormValue("points") == "" {
		points = -1
	} else {
		points, err = strconv.Atoi(r.FormValue("points"))
		if err != nil {
			h.handleError(w, r, "LocationEditPost: converting points", "Error converting points", "error", err)
			return
		}
	}

	// These are out of range values that will be ignored
	lat, lng := 200.0, 200.0
	if r.FormValue("latitude") != "" {
		lat, err = strconv.ParseFloat(r.FormValue("latitude"), 64)
		if err != nil {
			h.handleError(w, r, "LocationEditPost: converting latitude", "Error converting latitude", "error", err)
			return
		}
	}

	if r.FormValue("longitude") != "" {
		lng, err = strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			h.handleError(w, r, "LocationEditPost: converting longitude", "Error converting longitude", "error", err)
			return
		}
	}

	if lat > 90 || lat < -90 || lng > 180 || lng < -180 && (lat != 200.0 && lng != 200.0) {
		h.handleError(w, r, "LocationEditPost: invalid coordinates", "Invalid coordinates", "error", nil)
		return
	}

	data := services.LocationUpdateData{
		Name:      r.FormValue("name"),
		Latitude:  lat,
		Longitude: lng,
		Points:    points,
	}

	location, err := h.LocationService.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, locationCode)
	if err != nil {
		h.handleError(w, r, "LocationEditPost: finding location", "Error finding location", "error", err)
		return
	}

	err = h.LocationService.UpdateLocation(r.Context(), location, data)
	if err != nil {
		h.handleError(w, r, "LocationEditPost: updating location", "Error updating location", "error", err)
		return
	}

	if user.CurrentInstance.Settings.NavigationMethod == models.ShowClues {
		// Fetch the form clues
		clues := []string{}
		for key, value := range r.Form {
			if key == "clues" {
				clues = value
			}
		}
		clueIDs := []string{}
		for key, value := range r.Form {
			if key == "clue-ids" {
				clueIDs = value
			}
		}

		err = h.ClueService.UpdateClues(r.Context(), location, clues, clueIDs)
		if err != nil {
			h.handleError(w, r, "LocationEdit: updating clues", "Error updating clues", "error", err, "instance_id", user.CurrentInstanceID, "location_id", location.ID)
			return
		}
	}

	h.handleSuccess(w, r, "Location updated")
}

// LocationDelete handles deleting a location
func (h *AdminHandler) LocationDelete(w http.ResponseWriter, r *http.Request) {
	locationCode := chi.URLParam(r, "id")

	user := h.UserFromContext(r.Context())
	user.CurrentInstance.LoadLocations(r.Context())

	// Make sure the location exists and is part of the current instance
	var location models.Location
	for _, l := range user.CurrentInstance.Locations {
		if l.MarkerID == locationCode {
			location = l
			break
		}
	}
	if location.MarkerID == "" {
		h.Logger.Error("LocationDelete: finding location", "error", "location not found")
		err := templates.Toast(*flash.NewError("Location not found")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("LocationDelete: rendering toast:", "error", err)
		}
		return
	}

	err := h.GameManagerService.DeleteLocation(r.Context(), &location)
	if err != nil {
		h.Logger.Error("LocationDelete: deleting location", "error", err)
		err := templates.Toast(*flash.NewError("Error deleting location")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("LocationDelete: rendering toast:", "error", err)
		}
		return
	}

	w.Header().Set("HX-Redirect", "/admin/locations")
}
