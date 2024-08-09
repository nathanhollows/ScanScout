package server

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/routes"
	"github.com/nathanhollows/Rapua/internal/services"
)

var router *chi.Mux
var server *http.Server

func Start() {
	gameplayService := &services.GameplayService{}
	gameManagerService := &services.GameManagerService{}
	notificationService := services.NewNotificationService()

	router = routes.SetupRouter(gameplayService, gameManagerService, notificationService)

	server = &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}
	slog.Info("Server started", "addr", os.Getenv("SERVER_ADDR"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
