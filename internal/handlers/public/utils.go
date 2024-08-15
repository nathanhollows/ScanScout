package handlers

import (
	"log/slog"
)

type PublicHandler struct {
	Logger *slog.Logger
}

func NewPublicHandler(logger *slog.Logger) *PublicHandler {
	return &PublicHandler{
		Logger: logger,
	}
}
