package handlers

import (
	"fmt"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/sessions"
	templates "github.com/nathanhollows/Rapua/internal/templates/blocks"
)

// ValidateBlock runs input validation on the block
func (h *PlayerHandler) ValidateBlock(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r, "scanscout")
	teamCode := session.Values["team"]

	// Attempt to fetch team if teamCode is present
	team, err := h.getTeamIfExists(r.Context(), teamCode)
	if err != nil || team == nil {
		// If the team is not found, return an error
		h.handleError(w, r, fmt.Errorf("validateBlock: getTeamifExists: %v", err).Error(), "Team not found")
		err := invalidateSession(r, w)
		if err != nil {
			h.handleError(w, r, fmt.Errorf("validateBlock: invalidateSession: %v", err).Error(), "Something went wrong!")
		}
		return
	}

	r.ParseForm()
	data := make(map[string][]string)
	for key, value := range r.Form {
		data[key] = value
	}

	state, block, err := h.GameplayService.ValidateAndUpdateBlockState(r.Context(), *team, data)
	if err != nil {
		h.Logger.Error("validateBlock: validating and updating block state", "Something went wrong. Please try again.", err, "block", block.GetID(), "team", team.Code)
		return
	}

	err = templates.RenderPlayerUpdate(team.Instance.Settings, block, state).Render(r.Context(), w)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: rendering template: %v", err).Error(), "Something went wrong!")
		return
	}
}
