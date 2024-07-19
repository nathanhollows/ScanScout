package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/handlers"
)

func (h *PublicHandler) Index(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)
	data := handlers.TemplateData(r)
	data["title"] = "Home"
	handlers.Render(w, data, handlers.PublicDir, "home")
}
