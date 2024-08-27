package handlers

import (
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/templates/public"
)

// NotFound shows the not found page
func (h *PublicHandler) Pricing(w http.ResponseWriter, r *http.Request) {
	c := templates.Pricing()
	err := templates.PublicLayout(c, "Pricing").Render(r.Context(), w)

	if err != nil {
		h.Logger.Error("rendering Pricing page", "err", err)
	}
}
