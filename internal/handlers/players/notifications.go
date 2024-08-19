package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// DismissNotificationPost dismisses a message
func (h *PlayerHandler) DismissNotificationPost(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	notificationID := chi.URLParam(r, "ID")
	err := h.NotificationService.DismissNotification(r.Context(), notificationID)

	// Handle HTMX request
	if r.Header.Get("HX-Request") == "true" {
		if err != nil {
			h.Logger.Error("dismissing notification", "error", err.Error(), "notificationID", notificationID)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if err != nil {
		h.Logger.Error("dismissing notification", "error", err.Error(), "notificationID", notificationID)
		flash.NewError("Error dismissing notification").Save(w, r)
		http.Redirect(w, r, r.Header.Get("referer"), http.StatusSeeOther)

		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
