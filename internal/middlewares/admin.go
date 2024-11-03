package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/flash"
	internalServices "github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/services"
)

// AdminAuthMiddleware ensures the user is authenticated and has verified their email.
func AdminAuthMiddleware(authService internalServices.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure the user is authenticated
		user, err := authService.GetAuthenticatedUser(r)
		if err != nil {
			flash.NewError("You must be logged in to access this page").Save(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Redirect to verify email if the user hasn't verified their email
		// and they didn't sign up with OAuth
		if !user.EmailVerified && user.Provider == "" {
			http.Redirect(w, r, "/verify-email", http.StatusSeeOther)
			return
		}

		// Add the user to the context
		ctx := context.WithValue(r.Context(), contextkeys.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminCheckInstanceMiddleware ensures the user has an instance selected.
func AdminCheckInstanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextkeys.UserKey).(*models.User)

		// Check if the route contains /admin/instances
		reg := regexp.MustCompile(`/admin/instances/?`)
		if reg.MatchString(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if user.CurrentInstanceID == "" {
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

// AdminCheckBillingMiddleware ensures the user has a billing plan selected.
func AdminCheckBillingMiddleware(billingService services.BillingService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(contextkeys.UserKey).(*models.User)

		required, err := billingService.RequiresPlanSelection(r.Context(), user.ID)
		if err != nil {
			slog.Error("Failed to check billing status", err)
			flash.NewError("Failed to check billing status").Save(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if required {
			// If url is /admin/payment, don't redirect
			if r.URL.Path == "/admin/payment" {
				next.ServeHTTP(w, r)
				return
			} else {
				http.Redirect(w, r, "/admin/payment", http.StatusSeeOther)
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}
