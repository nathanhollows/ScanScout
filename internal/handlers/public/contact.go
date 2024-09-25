package handlers

import (
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/templates/public"
)

func (h *PublicHandler) Contact(w http.ResponseWriter, r *http.Request) {
	c := templates.Contact()
	err := templates.PublicLayout(c, "Contact").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Contact: rendering template", "error", err)
	}
}
