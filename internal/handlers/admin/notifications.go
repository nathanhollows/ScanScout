package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/handlers"
)

// NotifyAllPost sends a notification to all teams
func (h *AdminHandler) NotifyAllPost(w http.ResponseWriter, r *http.Request) {
	handlers.SetDefaultHeaders(w)

	if err := r.ParseForm(); err != nil {
		flash.NewError("Error parsing form").Save(w, r)
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
		return
	}

	user := h.UserFromContext(r.Context())
	user.CurrentInstance.LoadTeams(r.Context())
	content := r.FormValue("content")

	// Send the notification
	err := h.NotificationService.SendNotificationToAll(r.Context(), user.CurrentInstance.Teams, content)
	if err != nil {
		slog.Error("sending notification", "err", err.Error())
		flash.NewError("Error sending announcement").Save(w, r)
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("Announcement sent").Save(w, r)
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}
