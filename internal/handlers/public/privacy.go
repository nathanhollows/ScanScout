package handlers

import (
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/templates/public"
)

func (h *PublicHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	c := templates.Privacy()
	err := templates.PublicLayout(c, "Privacy").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Error rendering index", "err", err)
		return
	}
}
