package handlers

import (
	"fmt"
	"net/http"

	templates "github.com/nathanhollows/Rapua/internal/blocks/templates"
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

	block, state, err := h.BlockService.GetBlockWithStateByBlockIDAndTeamCode(r.Context(), data["block"], team.Code)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: getting block with state: %v", err).Error(), "Something went wrong!")
		return
	}

	if state.IsComplete {
		h.handleSuccess(w, r, "Block already completed")
	}

	if state.BlockID == "" {
		state.BlockID = block.GetID()
		state.TeamCode = team.Code
		state.PlayerData = []byte("{}")
	}

	err = h.GameplayService.ValidateAndUpdateBlockState(r.Context(), block, &state, data)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: validating block: %v", err).Error(), "Something went wrong!")
		return
	}

	err = templates.RenderPlayerView(team.Instance.Settings, block, state).Render(r.Context(), w)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: rendering template: %v", err).Error(), "Something went wrong!")
		return
	}

	return
}
