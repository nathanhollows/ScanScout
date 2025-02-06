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
		h.handleError(w, r, "Instances: rendering template", "Error rendering template", "error", err, "instance_id", user.CurrentInstanceID)
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
	instance, err := h.IntanceService.CreateInstance(r.Context(), name, user)
	if err != nil {
		h.handleError(w, r, "InstancesCreate: creating instance", "Error creating instance", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	// Switch to the new instance
	_, err = h.IntanceService.SwitchInstance(r.Context(), user, instance.ID)
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

	instance, err := h.IntanceService.DuplicateInstance(r.Context(), user, id, name)
	if err != nil {
		h.handleError(w, r, "InstanceDuplicate: duplicating instance", "Error duplicating instance", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	_, err = h.IntanceService.SwitchInstance(r.Context(), user, instance.ID)
	if err != nil {
		h.handleError(w, r, "InstanceDuplicate: switching instance", "Error switching instance", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	h.redirect(w, r, "/admin/instances")
}

// InstanceSwitch switches the current instance
func (h *AdminHandler) InstanceSwitch(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	instanceID := chi.URLParam(r, "id")
	if instanceID == "" {
		h.handleError(w, r, "InstanceSwitch: missing instance ID", "Could not switch instance", "error", "Instance ID is required", "instance_id", user.CurrentInstanceID)
		return
	}

	_, err := h.IntanceService.SwitchInstance(r.Context(), user, instanceID)
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

	_, err := h.IntanceService.DeleteInstance(r.Context(), user, id, confirmName)
	if err != nil {
		h.handleError(w, r, "InstanceDelete: deleting instance", "Error deleting instance", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	h.redirect(w, r, "/admin/instances")
}
