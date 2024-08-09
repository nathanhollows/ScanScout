package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/handlers"
)

// NotFound shows the not found page
func (h *PlayerHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	data["title"] = "Not Found"

	handlers.Render(w, data, handlers.PlayerDir, "not_found")
}
