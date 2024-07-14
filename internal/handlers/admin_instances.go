package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

var gameManagerService = services.NewGameManagerService()

// AdminInstancesHandler shows admin the instances
func AdminInstancesHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Instances"

	data["messages"] = flash.Get(w, r)
	render(w, data, true, "instances_index")
}

// AdminInstanceCreateHandler creates a new instance
func AdminInstanceCreateHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "New Instance"

	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.NewError("User not authenticated").Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	_, err := gameManagerService.CreateInstance(r.Context(), name, user)
	if err != nil {
		flash.NewError("Error creating instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Instance created successfully").Save(w, r)
	http.Redirect(w, r, "/admin/instances/", http.StatusSeeOther)
}

// AdminInstanceSwitchHandler switches the current instance
func AdminInstanceSwitchHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Switch Instance"

	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.NewError("User not authenticated").Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		flash.NewError("Instance ID is required").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	instance, err := gameManagerService.SwitchInstance(r.Context(), user, id)
	if err != nil {
		flash.NewError("Error switching instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("You are now using "+instance.Name+" as your current instance").Save(w, r)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// AdminInstanceDeleteHandler deletes an instance
func AdminInstanceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "Delete Instance"

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	id := r.Form.Get("id")
	confirmName := r.Form.Get("name")

	user, ok := data["user"].(*models.User)
	if !ok || user == nil {
		flash.NewError("User not authenticated").Save(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := gameManagerService.DeleteInstance(r.Context(), user, id, confirmName)
	if err != nil {
		flash.NewError("Error deleting instance: "+err.Error()).Save(w, r)
		http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Instance deleted successfully").Save(w, r)
	http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
}
