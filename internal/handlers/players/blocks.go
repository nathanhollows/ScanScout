package handlers

import (
	"fmt"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/sessions"
)

// ValidateBlock runs input validation on the block
func (h *PlayerHandler) ValidateBlock(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r, "scanscout")
	teamCode := session.Values["team"]

	// Attempt to fetch team if teamCode is present
	team, err := h.getTeamIfExists(r, teamCode)
	if err != nil || team == nil {
		// If the team is not found, return an error
		h.handleError(w, r, fmt.Errorf("validateBlock: getTeamifExists: %v", err).Error(), "Team not found")
		invalidateSession(session, r, w)
		return
	}

	r.ParseForm()
	data := make(map[string]string)
	for key, value := range r.Form {
		data[key] = value[0]
	}

	return
}
