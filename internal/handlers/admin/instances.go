package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// Instances shows admin the instances
func (h *AdminHandler) Instances(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Instances"
	data["page"] = "instances"

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "instances_index")
}

// InstancesCreate creates a new instance
func (h *AdminHandler) InstancesCreate(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "New Instance"

	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	_, err := h.GameManagerService.CreateInstance(r.Context(), name, user)
	if err != nil {
		flash.NewError("Error creating instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Instance created successfully").Save(w, r)
	http.Redirect(w, r, "/admin/instances/", http.StatusSeeOther)
}

// InstanceDuplicate duplicates an instance
func (h *AdminHandler) InstanceDuplicate(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	r.ParseForm()

	id := r.Form.Get("id")
	name := r.Form.Get("name")

	response := h.GameManagerService.DuplicateInstance(r.Context(), user, id, name)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		slog.Error("duplicating instance", "error", response.Error.Error())
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)
		return
	}

	newInstanceID := response.Data["instanceID"].(string)
	_, err := h.GameManagerService.SwitchInstance(r.Context(), user, newInstanceID)
	if err != nil {
		flash.NewError("Error switching instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
}

// InstanceSwitch switches the current instance
func (h *AdminHandler) InstanceSwitch(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Switch Instance"

	user := h.UserFromContext(r.Context())

	id := chi.URLParam(r, "id")
	if id == "" {
		flash.NewError("Instance ID is required").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	instance, err := h.GameManagerService.SwitchInstance(r.Context(), user, id)
	if err != nil {
		flash.NewError("Error switching instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("You are now using "+instance.Name+" as your current instance").Save(w, r)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// InstanceDelete deletes an instance
func (h *AdminHandler) InstanceDelete(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Delete Instance"

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	id := r.Form.Get("id")
	confirmName := r.Form.Get("name")

	user := h.UserFromContext(r.Context())

	err := h.GameManagerService.DeleteInstance(r.Context(), user, id, confirmName)
	if err != nil {
		flash.NewError("Error deleting instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Instance deleted successfully").Save(w, r)
	http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
}
