package handlers

import (
	"net/http"

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

	// Get the list of instances
	instances, err := models.FindAllInstances(user.UserID)
	if err != nil {
		flash.Message{
			Title:   "Error",
			Message: err.Error(),
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	data["instances"] = instances

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
