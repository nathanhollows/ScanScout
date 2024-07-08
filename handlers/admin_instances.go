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

	// Get the list of instances
	instances, err := models.FindAllInstances(data["user"].(*models.User).UserID)
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
	render(w, data, true, "instances_index")
}
