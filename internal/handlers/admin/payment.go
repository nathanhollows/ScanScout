package handlers

import (
	"net/http"

	"github.com/nathanhollows/Rapua/internal/templates/admin"
	p "github.com/nathanhollows/Rapua/internal/templates/public"
)

func (h *AdminHandler) Payment(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	c := templates.Payment(user.Email)
	err := p.AuthLayout(c, "Payment").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Payment: rendering template", "error", err)
	}
}

func (h *AdminHandler) PaymentPost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	r.ParseForm()

	c := templates.Payment(user.Email)
	err := p.AuthLayout(c, "Payment").Render(r.Context(), w)
	if err != nil {
		h.Logger.Error("Payment: rendering template", "error", err)
	}
}
