package handlers

import (
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/templates/public"
	"github.com/nathanhollows/Rapua/services"
)

func (h *PublicHandler) Docs(w http.ResponseWriter, r *http.Request) {
	docsService, err := services.NewDocsService("./docs")
	if err != nil {
		h.Logger.Error("Docs: creating docs service", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract the path after /docs/
	path := r.URL.Path
	if path == "/docs" || path == "/docs/" {
		path = "/docs/index"
	}

	// Get the page from the DocsService
	page, err := docsService.GetPage(path)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	c := templates.Docs(page, docsService.Pages)
	err = templates.PublicLayout(c, page.Title+" - Docs").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Contact: rendering template", "error", err)
	}
}
