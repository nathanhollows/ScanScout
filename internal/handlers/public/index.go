package handlers

import (
	"net/http"

	templates "github.com/nathanhollows/Rapua/v3/internal/templates/public"
)

func (h *PublicHandler) Index(w http.ResponseWriter, r *http.Request) {
	c := templates.Index()
	err := templates.PublicLayout(c, "Home").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Error rendering index", "err", err)
		return
	}
}
