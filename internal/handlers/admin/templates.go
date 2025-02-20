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

	_, err := h.TemplateService.CreateFromInstance(r.Context(), user.ID, user.CurrentInstanceID, name)
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
