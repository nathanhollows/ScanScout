package handlers

import (
	"net/http"
)

// NotifyAllPost sends a notification to all teams
func (h *AdminHandler) NotifyAllPost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "NotifyAllPost parsing form", "Error parsing form", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	content := r.FormValue("content")

	// Send the notification
	err := h.NotificationService.SendNotificationToAllTeams(r.Context(), user.CurrentInstanceID, content)
	if err != nil {
		h.handleError(w, r, "NotifyAllPost sending notification", "Error sending notification", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	h.handleSuccess(w, r, "Notification sent")
}

// NotifyTeamPost sends a notification to a specific team
func (h *AdminHandler) NotifyTeamPost(w http.ResponseWriter, r *http.Request) {
	user := h.UserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.handleError(w, r, "NotifyTeamPost parsing form", "Error parsing form", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	content := r.FormValue("content")
	teamCode := r.FormValue("teamCode")

	// Send the notification
	_, err := h.NotificationService.SendNotification(r.Context(), teamCode, content)
	if err != nil {
		h.handleError(w, r, "NotifyTeamPost sending notification", "Error sending notification", "error", err, "instance_id", user.CurrentInstanceID)
		return
	}

	h.handleSuccess(w, r, "Notification sent")
}
