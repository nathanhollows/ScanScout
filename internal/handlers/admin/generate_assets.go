package handlers

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
)

// QRCode handles the generation of QR codes for the current instance.
func (h *AdminHandler) QRCode(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextkeys.UserKey).(*models.User)

	// Extract parameters from the URL
	extension := chi.URLParam(r, "extension")
	if extension != "png" && extension != "svg" {
		h.Logger.Error("QRCodeHandler: Invalid extension provided")
		http.Error(w, "Invalid extension provided", http.StatusNotFound)
		return
	}

	action := chi.URLParam(r, "action")
	if action != "in" && action != "out" {
		h.Logger.Error("QRCodeHandler: Invalid type provided")
		http.Error(w, "Improper type provided", http.StatusNotFound)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("QRCodeHandler: No location provided")
		http.Error(w, "No location provided", http.StatusNotFound)
		return
	}

	// Check if the location exists
	if !h.GameManagerService.ValidateLocationMarker(user, id) {
		h.Logger.Error("QRCodeHandler: Location not found", "location", id)
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	// Get the path and content for the QR code
	path, content := h.GameManagerService.GetQRCodePathAndContent(action, id, "", extension)

	// Check if the file already exists, if so serve it
	if _, err := os.Stat(path); err == nil {
		if extension == "svg" {
			w.Header().Set("Content-Type", "image/svg+xml")
		} else {
			w.Header().Set("Content-Type", "image/png")
		}
		http.ServeFile(w, r, path)
		return
	}

	// Generate the QR code
	err := h.AssetGenerator.CreateQRCodeImage(
		r.Context(),
		path,
		content,
		h.AssetGenerator.WithFormat(extension),
	)
	if err != nil {
		h.Logger.Error("QRCodeHandler: Could not create QR code", "error", err)
		http.Error(w, "Could not create QR code", http.StatusInternalServerError)
		return
	}

	// Serve the generated QR code
	switch extension {
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case "png":
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "Invalid extension provided", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)

}

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

	var paths []string
	actions := []string{"in"}
	if user.CurrentInstance.Settings.CompletionMethod == models.CheckInAndOut {
		actions = []string{"in", "out"}
	}
	for _, location := range user.CurrentInstance.Locations {
		for _, extension := range []string{"png", "svg"} {
			for _, action := range actions {
				path, content := h.GameManagerService.GetQRCodePathAndContent(action, location.MarkerID, location.Name, extension)
				paths = append(paths, path)

				// Check if the file already exists, otherwise generate it
				if _, err := os.Stat(path); err == nil {
					continue
				}

				// Generate the QR code
				err := h.AssetGenerator.CreateQRCodeImage(
					r.Context(),
					path,
					content,
					h.AssetGenerator.WithFormat(extension),
				)
				if err != nil {
					h.Logger.Error("QRCodeHandler: Could not create QR code", "error", err)
					http.Error(w, "Could not create QR code", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	path, err := h.AssetGenerator.CreateArchive(r.Context(), paths)
	if err != nil {
		h.Logger.Error("QR codes could not be zipped", "error", err, "instance", user.CurrentInstanceID)
		http.Error(w, "QR codes could not be zipped", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, path)
	os.Remove(path)
}
