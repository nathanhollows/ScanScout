package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/helpers"
	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// Instances shows admin the instances
func (h *AdminHandler) MarkdownGuide(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.MarkdownGuide()
	err := templates.Layout(c, *user, "Markdown", "Markdown Guide").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("MarkdownGuide: rendering template", "error", err)
	}
}

// PreviewMarkdown takes markdown from a form and renders it for htmx
func (h *AdminHandler) PreviewMarkdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var m map[string]string
	err := decoder.Decode(&m)
	if err != nil {
		h.Logger.Error("markdown preview: decoding JSON", "error", err)
		err := templates.Toast(*flash.NewError("Error converting markdown")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("markdown preview: rendering template", "error", err)
		}
		return
	}

	md, err := helpers.MarkdownToHTML(m["markdown"])
	if err != nil {
		h.Logger.Error("markdown preview: converting string to markdown", "error", err)
		err := templates.Toast(*flash.NewError("Error converting markdown")).Render(r.Context(), w)
		if err != nil {
			h.Logger.Error("markdown preview: rendering template", "error", err)
		}
		return
	}

	err = templates.MarkdownPreview(md).Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("markdown preview: rendering template", "error", err)
	}

}
