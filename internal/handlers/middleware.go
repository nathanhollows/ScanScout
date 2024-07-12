package handlers

import (
	"context"
	"net/http"
	"regexp"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
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

		ctx := context.WithValue(r.Context(), models.UserIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Ensure the user has an instance selected, otherwise redirect to the instances page
func adminCheckInstanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.UserIDKey).(*models.User)

		// Check if the route contains /admin/instances
		reg := regexp.MustCompile(`/admin/instances/?`)
		if reg.MatchString(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if user.CurrentInstance == nil {
			flash.Message{
				Title:   "Error",
				Message: "Please select an instance to continue",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/admin/instances", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
