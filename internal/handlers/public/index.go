package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/templates"
)

func (h *PublicHandler) Index(w http.ResponseWriter, r *http.Request) {
	c := templates.Index()
	err := templates.PublicLayout(c, "Home").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Error rendering index", "err", err)
		return
	}
}
