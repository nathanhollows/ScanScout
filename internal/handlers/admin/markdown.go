package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// Instances shows admin the instances
func (h *AdminHandler) MarkdownGuide(w http.ResponseWriter, r *http.Request) {
	data := handlers.TemplateData(r)
	data["title"] = "Markdown Guide"
	data["page"] = "markdown"

	data["messages"] = flash.Get(w, r)
	handlers.Render(w, data, handlers.AdminDir, "markdown")
}

// PreviewMarkdown takes markdown from a form and renders it for htmx
func (h *AdminHandler) PreviewMarkdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)

	decoder := json.NewDecoder(r.Body)
	var m map[string]string
	err := decoder.Decode(&m)
	if err != nil {
		h.Logger.Error("markdown preview: decoding JSON", "error", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	data["markdown"] = m["markdown"]

	handlers.RenderHTMX(w, data, handlers.AdminDir, "markdown_preview")
}
