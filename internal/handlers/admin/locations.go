package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Locations shows admin the locations
func (h *AdminHandler) Locations(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Locations"
	data["page"] = "locations"

	user := h.UserFromContext(r.Context())
	user.CurrentInstance.Locations.LoadClues(r.Context())
	data["locations"] = user.CurrentInstance.Locations

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "locations_index")
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
	err := r.ParseForm()
	if err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/locations/new", http.StatusSeeOther)
		return
	}

	user := h.UserFromContext(r.Context())

	name := r.FormValue("name")
	content := r.FormValue("content")
	criteriaID := r.FormValue("criteria")
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")

	response := h.GameManagerService.CreateLocation(r.Context(), user, name, content, criteriaID, lat, lng)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("creating location", "error", response.Error.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	location, ok := response.Data["location"].(*models.Location)
	if !ok {
		h.Logger.Error("creating location", "error", "location not returned")
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/locations/"+location.MarkerID, http.StatusSeeOther)
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

	err := r.ParseForm()
	if err != nil {
		h.Logger.Error("reordering locations", "error", err.Error())
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	user := h.UserFromContext(r.Context())

	locations := r.Form["location"]
	response := h.GameManagerService.ReorderLocations(r.Context(), user, locations)
	// Discard the flash messages since this is invoked via HTMX
	if response.Error != nil {
		h.Logger.Error("reordering locations", "error", response.Error.Error())
		http.Error(w, response.Error.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, "Reordered locations", http.StatusOK)
}

// LocationEdit shows the form to edit a location
func (h *AdminHandler) LocationEdit(w http.ResponseWriter, r *http.Request) {
	// Get the location from the chi context
	code := chi.URLParam(r, "id")

	user := h.UserFromContext(r.Context())

	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, code)
	if err != nil {
		h.Logger.Error("finding location", "error", err.Error())
		flash.NewError("Location could not be found").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	location.LoadClues(r.Context())

	c := templates.EditLocation(*location, user.CurrentInstance.Settings)
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

	err = h.GameManagerService.UpdateLocation(r.Context(), location, newName, newContent, lat, lng)
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
