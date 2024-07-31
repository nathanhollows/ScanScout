package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// Instances shows admin the instances
func (h *AdminHandler) MarkdownGuide(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Markdown Guide"
	data["page"] = "markdown"

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "markdown")
}
