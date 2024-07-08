package handlers

import (
	"context"
	"net/http"

	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Get(r, "admin")
		if err != nil {
			http.Error(w, "Session error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if session.Values["user_id"] == nil {
			flash.Message{
				Title:   "Error",
				Message: "You must be logged in to access this page",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := models.FindUserBySession(r)
		if err != nil {
			http.Error(w, "User not found: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
