package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
)

func adminGenerateQRHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "id")

	user := r.Context().Value(contextkeys.UserKey).(*models.User)

	location, err := models.FindLocationByInstanceAndCode(r.Context(), user.CurrentInstanceID, code)
	if err != nil {
		flash.NewError("Location could not be found").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	err = location.Marker.GenerateQRCode()
	if err != nil {
		flash.NewError("QR code could not be generated").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}
}

func AdminLocationQRZipHandler(w http.ResponseWriter, r *http.Request) {
	archive, err := models.GenerateQRCodeArchive(r.Context(), r.Context().Value(contextkeys.UserKey).(*models.User).CurrentInstanceID)
	if err != nil {
		flash.NewError("QR codes could not be zipped").Save(w, r)
		http.Redirect(w, r, "/admin/locations", http.StatusSeeOther)
		return
	}

	// Overwrite the file download header
	instance := r.Context().Value(contextkeys.UserKey).(*models.User).CurrentInstance
	w.Header().Set("Content-Disposition", "attachment; filename="+instance.Name+" QR codes .zip")
	// Serve the file
	http.ServeFile(w, r, archive)

}
