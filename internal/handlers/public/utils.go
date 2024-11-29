package handlers

import (
	"log/slog"

	"github.com/nathanhollows/Rapua/internal/services"
)

type PublicHandler struct {
	Logger      *slog.Logger
	AuthService services.AuthService
	UserService services.UserService
}

func NewPublicHandler(
	logger *slog.Logger,
	authService services.AuthService,
	userService services.UserService,
) *PublicHandler {
	return &PublicHandler{
		Logger:      logger,
		AuthService: authService,
		UserService: userService,
	}
}
