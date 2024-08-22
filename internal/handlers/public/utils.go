package handlers

import (
	"log/slog"

	"github.com/nathanhollows/Rapua/internal/services"
)

type PublicHandler struct {
	Logger       *slog.Logger
	UserServices services.UserServices
}

func NewPublicHandler(logger *slog.Logger, userServices services.UserServices) *PublicHandler {
	return &PublicHandler{
		Logger:       logger,
		UserServices: userServices,
	}
}
