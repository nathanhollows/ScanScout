package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/v3/internal/flash"
	templates "github.com/nathanhollows/Rapua/v3/internal/templates/admin"
)

// TemplatesCreate creates a new template, which is a type of instance.
func (h *AdminHandler) TemplatesCreate(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "TemplateCreate: parsing form", "Error parsing form", "error", err)
		return
	}

	name := r.FormValue("name")
	id := r.FormValue("id")

	if id == "" {
		h.handleError(w, r, "TemplateCreate: missing id", "Could not find the instance ID")
		return
	}
	if name == "" {
		h.handleError(w, r, "TemplateCreate: missing name", "Please provide a name for the template")
		return
	}

	_, err := h.TemplateService.CreateFromInstance(r.Context(), user.ID, id, name)
	if err != nil {
		h.handleError(w, r, "TemplateCreate: creating instance", "Error creating instance", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	err = templates.Toast(*flash.NewSuccess("Template created")).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("InstanceDelete: rendering template", "Error", err)
	}

	gameTemplates, err := h.TemplateService.Find(r.Context(), user.ID)
	err = templates.Templates(gameTemplates).Render(r.Context(), w)
	if err != nil {
		h.handleError(w, r, "Instances: rendering template", "Error rendering template", "error", err, "instance_id", user.CurrentInstanceID)
	}
}

// TemplatesLaunch launches an instance from a template.
func (h *AdminHandler) TemplatesLaunch(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "TemplateCreate: parsing form", "Error parsing form", "error", err)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		h.handleError(w, r, "TemplateCreate: missing id", "Could not find the instance ID")
		return
	}

	name := r.FormValue("name")
	if name == "" {
		h.handleError(w, r, "TemplateCreate: missing name", "Please provide a name for the template")
		return
	}

	// Regenerate refers to location codes
	regen := r.Form.Has("regenerate")

	// Create a new instance from the template
	newGame, err := h.TemplateService.LaunchInstance(r.Context(), user.ID, id, name, regen)
	if err != nil {
		h.handleError(w, r, "TemplateCreate: creating instance", "Error creating instance", "error", err, "user_id", user.ID)
		return
	}

	// Switch to the new instance
	_, err = h.InstanceService.SwitchInstance(r.Context(), user, newGame.ID)
	if err != nil {
		h.handleError(w, r, "InstancesCreate: switching instance", "Error switching instance", "error", err)
		return
	}

	h.redirect(w, r, "/admin/instances")
}

// TemplatesDelete deletes a template.
func (h *AdminHandler) TemplatesDelete(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "TemplateDelete: parsing form", "Error parsing form", "error", err)
		return
	}

	id := r.Form.Get("id")
	if id == "" {
		h.handleError(w, r, "TemplateDelete: missing id", "Could not find the instance ID")
		return
	}

	template, err := h.TemplateService.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, "TemplateDelete: getting template", "Error getting template", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	_, err = h.InstanceService.DeleteInstance(r.Context(), user, template.ID, template.Name)
	if err != nil {
		h.handleError(w, r, "InstanceDelete: deleting instance", "Error deleting instance", "error", err, "instance_id", user.CurrentInstanceID)
	} else {
		err = templates.Toast(*flash.NewSuccess("Template deleted")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("InstanceDelete: rendering template", "Error", err)
		}
	}

	gameTemplates, err := h.TemplateService.Find(r.Context(), user.ID)
	err = templates.Templates(gameTemplates).Render(r.Context(), w)
	if err != nil {
		h.handleError(w, r, "Instances: rendering template", "Error rendering template", "error", err, "instance_id", user.CurrentInstanceID)
	}
}
