package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Locations shows admin the locations
func (h *AdminHandler) Locations(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.LocationsIndex(user.CurrentInstance.Settings, user.CurrentInstance.Locations)
	err := templates.Layout(c, *user, "Locations", "Locations").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Locations: rendering template", "error", err)
	}

}

// LocationNew shows the form to create a new location
func (h *AdminHandler) LocationNew(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.AddLocation()
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

	response := h.GameManagerService.CreateLocation(r.Context(), user, data)
	if response.Error != nil {
		h.handleError(w, r, "LocationNewPost: creating location", "Error creating location", "error", response.Error, "instance_id", user.CurrentInstanceID)
		return
	}

	location, ok := response.Data["location"].(*models.Location)
	if !ok {
		h.handleError(w, r, "LocationNewPost: creating location", "Error creating location", "error", "location not returned", "instance_id", user.CurrentInstanceID)
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
	// Get the location from the chi context
	code := chi.URLParam(r, "id")

	user := h.UserFromContext(r.Context())

	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, code)
	if err != nil {
		h.handleError(w, r, "LocationEdit: finding location", "Error finding location", "error", err, "instance_id", user.CurrentInstanceID, "location_code", code)
		return
	}

	blocks, err := h.BlockService.GetByLocationID(r.Context(), location.ID)
	if err != nil {
		h.Logger.Error("LocationEdit: getting blocks", "error", err, "instance_id", user.CurrentInstanceID, "location_id", location.ID)
		h.redirect(w, r, "/admin/locations")
		return
	}

	location.LoadClues(r.Context())

	c := templates.EditLocation(*location, user.CurrentInstance.Settings, blocks)
	err = templates.Layout(c, *user, "Locations", "Edit Location").Render(r.Context(), w)
}

// LocationEditPost handles updating a location
func (h *AdminHandler) LocationEditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.Logger.Error("LocationEditPost: parsing form", "error", err)
		err := templates.Toast(*flash.NewError("Error parsing form")).Render(r.Context(), w)

		if err != nil {
			h.Logger.Error("LocaiotnEditPost: rendering toast:", "error", err)
		}
		return
	}

	locationCode := chi.URLParam(r, "id")

	user := h.UserFromContext(r.Context())

	location, err := models.FindLocationByInstanceAndCode(
		r.Context(),
		user.CurrentInstanceID,
		locationCode,
	)
	if err != nil {
		h.Logger.Error("LocationEditPost: finding location", "err", err)
		err := templates.Toast(*flash.NewError("Location not found")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("LocationEditPost: rendering toast:", "error", err)
		}
		return
	}

	newName := r.FormValue("name")
	newContent := r.FormValue("content")
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")
	pts := r.FormValue("points")
	points, err := strconv.Atoi(pts)
	if err != nil {
		h.Logger.Error("LocationEditPost: converting points", "error", err)
		err := templates.Toast(*flash.NewError("Error saving location")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("LocationEditPost: rendering toast:", "error", err)
		}
		return
	}

	err = h.GameManagerService.UpdateLocation(r.Context(), location, newName, newContent, lat, lng, points)
	if err != nil {
		h.Logger.Error("LocationEditPost: updating location", "error", err)
		err := templates.Toast(*flash.NewError("Error saving location")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("LocationEditPost: rendering toast:", "error", err)
		}
		return
	}

	err = h.GameManagerService.UpdateClues(r.Context(), location, r.Form["clues[]"], r.Form["clue_ids[]"])
	if err != nil {
		h.Logger.Error("LocationEditPost: updating clues", "error", err)
		err := templates.Toast(*flash.NewError("Error saving clues")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("LocationEditPost: rendering toast:", "error", err)
		}
		return
	}

	err = templates.Toast(*flash.NewSuccess("Location updated!")).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("LocationEditPost: rendering toast:", "error", err)
	}
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
