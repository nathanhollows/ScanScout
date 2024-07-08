package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
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
	// Instances are already loaded in the template data

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
		UserID: user.UserID,
	}

	if err := instance.Save(); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error saving instance",
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
	http.Redirect(w, r, "/admin/instances/"+instance.ID, http.StatusSeeOther)
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
	instance, err := models.FindInstanceByID(id)
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
	if err := user.Update(); err != nil {
		flash.Message{
			Title:   "Error",
			Message: "Error updating user",
			Style:   flash.Error,
		}.Save(w, r)
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
