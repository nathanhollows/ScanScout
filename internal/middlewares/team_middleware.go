package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// TeamMiddleware extracts the team code from the session and finds the matching instance.
func TeamMiddleware(teamService services.TeamService, next http.Handler) http.Handler {
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
		team, err := teamService.FindTeamByCode(r.Context(), teamCode)
		if err != nil {
			slog.Error("finding team by code: ", "err", err, "teamCode", teamCode)
			next.ServeHTTP(w, r)
			return
		}

		err = teamService.LoadRelation(r.Context(), team, "Instance")
		if err != nil {
			slog.Error("loading relations: ", "err", err)
			next.ServeHTTP(w, r)
			return
		}

		// Add team to context
		ctx := context.WithValue(r.Context(), contextkeys.TeamKey, team)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
