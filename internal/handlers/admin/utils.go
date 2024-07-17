package handlers

import "github.com/nathanhollows/Rapua/internal/services"

type AdminHandler struct {
	GameManagerService *services.GameManagerService
}

func NewAdminHandler(gs *services.GameManagerService) *AdminHandler {
	return &AdminHandler{
		GameManagerService: gs,
	}
}
