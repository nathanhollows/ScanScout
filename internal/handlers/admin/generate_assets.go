package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
)

// Show the form to edit the navigation settings.
func (h *AdminHandler) GeneratePosters(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextkeys.UserKey).(*models.User)

	posters, err := models.GeneratePosters(r.Context(), user.CurrentInstanceID)
	if err != nil {
		h.Logger.Error("Posters could not be generated", "error", err, "instance", user.CurrentInstanceID)

		flash.NewError("Posters could not be generated").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	// Overwrite the file download header
	instance := r.Context().Value(contextkeys.UserKey).(*models.User).CurrentInstance
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" posters.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	// Serve the file
	http.ServeFile(w, r, posters)
}
