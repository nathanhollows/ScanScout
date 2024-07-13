package middlewares

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// TeamMiddleware extracts the team code from the session and finds the matching instance.
func TeamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the session
		session, err := sessions.Get(r, "scanscout")
		if err != nil {
			log.Error("Error getting session: ", err)
			next.ServeHTTP(w, r)
			return
		}

		// Extract team code from session
		teamCode, ok := session.Values["team"].(string)
		if !ok || teamCode == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Find the matching team instance
		team, err := models.FindTeamByCode(r.Context(), teamCode)
		if err != nil {
			log.Error("Error finding team by code: ", err)
			next.ServeHTTP(w, r)
			return
		}

		// Add team to context
		ctx := context.WithValue(r.Context(), contextkeys.TeamKey, team)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
