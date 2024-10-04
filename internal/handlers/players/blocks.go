package handlers

import (
	"fmt"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/repositories"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// ValidateBlock runs input validation on the block
func (h *PlayerHandler) ValidateBlock(w http.ResponseWriter, r *http.Request) {
	blockStateRepo := repositories.NewBlockStateRepository()
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

	block, err := h.BlockService.GetByBlockID(r.Context(), data["block"])
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: getting block: %v", err).Error(), "Something went wrong!")
		return
	}

	state, err := blockStateRepo.GetByBlockAndTeam(r.Context(), block.GetID(), team.Code)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			h.handleError(w, r, "validateBlock: getting team block state", "Something went wrong!", "error", err)
			return
		}
	}

	if state.IsComplete {
		h.handleSuccess(w, r, "Block already completed")
	}

	err = h.GameplayService.ValidateAndUpdateBlockState(r.Context(), block, &state, data)
	if err != nil {
		h.handleError(w, r, fmt.Errorf("validateBlock: validating block: %v", err).Error(), "Something went wrong!")
		return
	}

	h.handleSuccess(w, r, "Block validated")

	return
}
