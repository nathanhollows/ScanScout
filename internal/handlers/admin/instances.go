package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
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
	data := handlers.TemplateData(r)
	data["title"] = "New Instance"

	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "InstancesCreate: parsing form", "Error parsing form", "error", err)
		return
	}

	name := r.FormValue("name")
	response := h.GameManagerService.CreateInstance(r.Context(), name, user)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
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
	data := handlers.TemplateData(r)
	data["title"] = "Switch Instance"

	user := h.UserFromContext(r.Context())

	id := chi.URLParam(r, "id")
	if id == "" {
		flash.NewError("Instance ID is required").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	_, err := h.GameManagerService.SwitchInstance(r.Context(), user, id)
	if err != nil {
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		err := templates.Toast(*flash.NewError("Error switching instance: " + err.Error())).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("InstanceSwitch: rendering template", "error", err)
		}
		return
	}

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

	if user.CurrentInstanceID == id {
		flash.NewError("You cannot delete the instance you are currently using").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	response := h.GameManagerService.DeleteInstance(r.Context(), user, id, confirmName)
	for _, message := range response.FlashMessages {
		message.Save(w, r)
	}
	if response.Error != nil {
		h.Logger.Error("deleting instance", "error", response.Error.Error())
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
}
