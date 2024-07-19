package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// Locations shows admin the locations
func (h *AdminHandler) Locations(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Locations"

	user := h.UserFromContext(r.Context())
	data["locations"] = user.CurrentInstance.Locations

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, true, "locations_index")
}

// LocationNew shows the form to create a new location
func (h *AdminHandler) LocationNew(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "New Location"
	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, true, "locations_new")
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

	http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
}
