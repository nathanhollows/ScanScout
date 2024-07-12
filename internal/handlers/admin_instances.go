package handlers

import (
	"log/slog"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
)

// adminInstancesHandler shows admin the instances
func adminInstancesHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Instances"

	// Get the current user
	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.Message{
			Title:   "Error",
			Message: "User not authenticated",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// NOTE:
	// data["instances"] are already loaded in the template data

	// Render the template
	data["messages"] = flash.Get(w, r)
	if err := render(w, data, true, "instances_index"); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// adminInstanceCreateHandler creates a new instance
func adminInstanceCreateHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "New Instance"

	// Get the current user
	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.Message{
			Title:   "Error",
			Message: "User not authenticated",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse the form
	if err := r.ParseForm(); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error parsing form",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		flash.Message{
			Title:   "Error",
			Message: "Name is required",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Create the new instance
	instance := &models.Instance{
		Name:   name,
		UserID: user.ID,
	}

	if err := instance.Save(r.Context()); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error saving instance",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Switch to the new instance
	user.CurrentInstanceID = instance.ID
	if err := user.Update(r.Context()); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error updating user",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Redirect to the new instance
	flash.Message{
		Title:   "Success",
		Message: "Instance created successfully",
		Style:   flash.Success,
	}.Save(w, r)
	http.Redirect(w, r, "/admin/instances/", http.StatusSeeOther)
}

// adminInstanceSwitchHandler switches the current instance
func adminInstanceSwitchHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Switch Instance"

	// Get the current user
	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.Message{
			Title:   "Error",
			Message: "User not authenticated",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get the instance ID from the URL
	id := chi.URLParam(r, "id")
	if id == "" {
		flash.Message{
			Title:   "Error",
			Message: "Instance ID is required",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Find the instance
	instance, err := models.FindInstanceByID(r.Context(), id)
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Instance not found",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Set the current instance
	user.CurrentInstanceID = instance.ID
	if err := user.Update(r.Context()); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error updating user",
			Style:   flash.Error,
		}.Save(w, r)
		slog.Error("Error updating user", "err", err.Error())
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Redirect to user back to the previous page
	flash.Message{
		Title:   "Success",
		Message: "You are now using " + instance.Name + " as your current instance",
		Style:   flash.Success,
	}.Save(w, r)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
	return

}

// adminInstanceDeleteHandler deletes an instance
func adminInstanceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Delete Instance"

	// Parse the form
	if err := r.ParseForm(); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error parsing form",
			Style:   flash.Error,
		}.Save(w, r)
		log.Error(err, "ctx", r.Context(), "form", r.Form)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Get the instance ID from the URL
	id := r.Form.Get("id")
	if id == "" {
		flash.Message{
			Title:   "Error",
			Message: "Instance ID is required",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Find the instance
	instance, err := models.FindInstanceByID(r.Context(), id)
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Instance not found",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Verify the user is the owner of the instance
	user := data["user"].(*models.User)
	if user.ID != instance.UserID {
		flash.Message{
			Title:   "Error",
			Message: "You do not have permission to delete this instance",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Before submitting the form, the user must confirm the deletion by
	// typing the instance name. If the name does not match, redirect back
	// to the form.
	if r.Form.Get("name") != instance.Name {
		flash.Message{
			Title:   "Error",
			Message: "Please type the instance name to confirm deletion",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances/"+instance.ID, http.StatusSeeOther)
		return
	}

	// Update the user's current instance if it matches the instance being deleted
	if user.CurrentInstanceID == instance.ID {
		user.CurrentInstanceID = ""
		if err := user.Update(r.Context()); err != nil {
			flash.Message{
				Title:   "Error",
				Message: "Error updating user",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
			return
		}
	}

	// Delete the instance
	if err := instance.Delete(r.Context()); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error deleting instance",
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	// Redirect to the instances page
	flash.Message{
		Title:   "Success",
		Message: "Instance deleted successfully",
		Style:   flash.Success,
	}.Save(w, r)
	http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
}
