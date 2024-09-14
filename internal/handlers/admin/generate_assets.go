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

	instance := user.CurrentInstance.Name
	w.Header().Set("Content-Disposition", "attachment; filename="+instance+" posters.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, posters)
}

// GenerateQRCodeArchive generates a zip file containing all the QR codes for the current instance.
func (h *AdminHandler) GenerateQRCodeArchive(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextkeys.UserKey).(*models.User)

	archive, err := models.GenerateQRCodeArchive(r.Context(), user.CurrentInstanceID)
	if err != nil {
		h.Logger.Error("QR codes could not be zipped", "error", err, "instance", user.CurrentInstanceID)
		flash.NewError("QR codes could not be zipped").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	instance := user.CurrentInstance.Name
	w.Header().Set("Content-Disposition", "attachment; filename="+instance+" QR codes .zip")
	http.ServeFile(w, r, archive)
}
