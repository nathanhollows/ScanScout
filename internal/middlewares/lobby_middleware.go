package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/models"
)

// LobbyMiddleware redirects to the lobby if the game is scheduled to start
func LobbyMiddleware(teamRepo repositories.TeamRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the session
		session, err := sessions.Get(r, "scanscout")
		if err != nil {
			slog.Error("getting session: ", "err", err, "ctx", r.Context())
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
		team, err := teamRepo.FindTeamByCode(r.Context(), teamCode)
		if err != nil {
			slog.Error("finding team by code: ", "err", err, "teamCode", teamCode)
			next.ServeHTTP(w, r)
			return
		}

		// Redirect to lobby if game is scheduled
		if team.Instance.GetStatus() != models.Active {
			http.Redirect(w, r, "/lobby", http.StatusFound)
			return
		}

		// Add team to context
		ctx := context.WithValue(r.Context(), contextkeys.TeamKey, team)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
