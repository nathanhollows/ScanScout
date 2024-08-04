package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Locations shows admin the locations
func (h *AdminHandler) Locations(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Locations"
	data["page"] = "locations"

	user := h.UserFromContext(r.Context())
	data["locations"] = user.CurrentInstance.Locations

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "locations_index")
}

// LocationNew shows the form to create a new location
func (h *AdminHandler) LocationNew(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "New Location"
	data["page"] = "locations"
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "locations_new")
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
		slog.Error("creating location", "error", response.Error.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	location, ok := response.Data["location"].(*models.Location)
	if !ok {
		slog.Error("creating location", "error", "location not returned")
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/locations/"+location.ID, http.StatusSeeOther)
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
		slog.Error("reordering locations", "error", err.Error())
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	user := h.UserFromContext(r.Context())

	locations := r.Form["location"]
	response := h.GameManagerService.ReorderLocations(r.Context(), user, locations)
	// Discard the flash messages since this is invoked via HTMX
	if response.Error != nil {
		slog.Error("reordering locations", "error", response.Error.Error())
		http.Error(w, response.Error.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, "Reordered locations", http.StatusOK)
}
