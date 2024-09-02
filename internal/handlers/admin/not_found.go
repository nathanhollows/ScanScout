package handlers

import (
	"log/slog"
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/templates/admin"
)

// NotFound shows the not found page
func (h *AdminHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.NotFound()
	err := templates.Layout(c, *user, "Error", "Not Found").Render(r.Context(), w)

	if err != nil {
		slog.Error("rendering NotFound page", "err", err)
	}
}
