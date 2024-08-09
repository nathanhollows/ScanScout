package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/handlers"
)

// NotFound shows the not found page
func (h *AdminHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Not Found"

	handlers.Render(w, data, handlers.AdminDir, "not_found")
}
