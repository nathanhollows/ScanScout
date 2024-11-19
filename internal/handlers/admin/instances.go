package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Instances shows admin the instances
func (h *AdminHandler) Instances(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.Instances(user.Instances, user.CurrentInstance)
	err := templates.Layout(c, *user, "Instances", "Instances").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Instances: rendering template", "error", err)
	}
}

// InstancesCreate creates a new instance
func (h *AdminHandler) InstancesCreate(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "InstancesCreate: parsing form", "Error parsing form", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	name := r.FormValue("name")
	response := h.GameManagerService.CreateInstance(r.Context(), name, user)
	if response.Error != nil {
		h.handleError(w, r, "InstancesCreate: creating instance", "Error creating instance", "error", response.Error)
		return
	}

	// Switch to the new instance
	_, err := h.GameManagerService.SwitchInstance(r.Context(), user, response.Data["instanceID"].(string))
	if err != nil {
		h.handleError(w, r, "InstancesCreate: switching instance", "Error switching instance", "error", err)
		return
	}

	h.redirect(w, r, "/admin/instances")
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
		h.Logger.Error("duplicating instance", "error", response.Error.Error())
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
	flash.NewSuccess("Now using "+name+" as your current instance").Save(w, r)

	http.Redirect(w, r, "/admin/navigation", http.StatusSeeOther)
}

// InstanceSwitch switches the current instance
func (h *AdminHandler) InstanceSwitch(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	id := chi.URLParam(r, "id")
	if id == "" {
		h.handleError(w, r, "InstanceSwitch: missing instance ID", "Could not switch instance", "error", "Instance ID is required", "instance_id", user.CurrentInstanceID)
		return
	}

	_, err := h.GameManagerService.SwitchInstance(r.Context(), user, id)
	if err != nil {
		h.handleError(w, r, "InstanceSwitch: switching instance", "Error switching instance", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	if r.URL.Query().Has("redirect") {
		h.redirect(w, r, r.URL.Query().Get("redirect"))
		return
	}

	h.redirect(w, r, r.Header.Get("Referer"))
}

// InstanceDelete deletes an instance
func (h *AdminHandler) InstanceDelete(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "InstanceDelete: parsing form", "Error parsing form", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	id := r.Form.Get("id")
	confirmName := r.Form.Get("name")

	if user.CurrentInstanceID == id {
		err := templates.Toast(*flash.NewError("You cannot delete the instance you are currently using")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("InstanceDelete: rendering template", "error", err)
		}
		return
	}

	response := h.GameManagerService.DeleteInstance(r.Context(), user, id, confirmName)
	if response.Error != nil {
		h.handleError(w, r, "InstanceDelete: deleting instance", "Error deleting instance", "error", response.Error, "instance_id", user.CurrentInstanceID)
		return
	}

	h.redirect(w, r, "/admin/instances")
}
