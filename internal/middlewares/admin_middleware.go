package middlewares

import (
	"context"
	"net/http"
	"regexp"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

// AdminAuthMiddleware ensures the user is authenticated.
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := services.GetAuthenticatedUser(r)
		if err != nil {
			flash.Message{
				Title:   "Error",
				Message: "You must be logged in to access this page",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminCheckInstanceMiddleware ensures the user has an instance selected.
func AdminCheckInstanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextkeys.UserIDKey).(*models.User)

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
