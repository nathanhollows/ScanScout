package handlers

import "github.com/nathanhollows/Rapua/internal/services"

type PlayerHandler struct {
	GameplayService *services.GameplayService
}

func NewPlayerHandler(gs *services.GameplayService) *PlayerHandler {
	return &PlayerHandler{
		GameplayService: gs,
	}
}
