package server

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/nathanhollows/Rapua/internal/routes"
	"golang.org/x/exp/slog"
)

var router *chi.Mux
var server *http.Server

func Start() {
	router = routes.SetupRouter()

	server = &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}
	slog.Info("Server started on %s", os.Getenv("SERVER_ADDR"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
