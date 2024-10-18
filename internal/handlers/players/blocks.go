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
	team, err := h.getTeamIfExists(r, teamCode)
	if err != nil || team == nil {
		// If the team is not found, return an error
		h.handleError(w, r, fmt.Errorf("validateBlock: getTeamifExists: %v", err).Error(), "Team not found")
		invalidateSession(session, r, w)
		return
	}

	r.ParseForm()
	data := make(map[string][]string)
	for key, value := range r.Form {
		data[key] = value
	}

	block, state, err := h.BlockService.GetBlockWithStateByBlockIDAndTeamCode(r.Context(), data["block"][0], team.Code)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: getting block with state: %v", err).Error(), "Something went wrong!")
		return
	}

	if state == nil {
		h.handleError(w, r, "validateBlock: getting block with state: state is nil", "Block state not found")
		return
	}

	if state.IsComplete() {
		h.handleSuccess(w, r, "Block already completed")
	}

	state, err = h.GameplayService.ValidateAndUpdateBlockState(r.Context(), block, state, data)
	if err != nil {
		h.Logger.Error("validateBlock: validating and updating block state", "error", err, "block", block.GetID(), "team", team.Code)
	}

	if state.IsComplete() {
		// err := h.GameplayService.UpdateCheckinStatus(r.Context(), team, state)
	}

	err = templates.RenderPlayerUpdate(team.Instance.Settings, block, state).Render(r.Context(), w)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: rendering template: %v", err).Error(), "Something went wrong!")
		return
	}

	return
}
