package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gomarkdown/markdown"
	"github.com/nathanhollows/ScanScout/models"
)

var router *chi.Mux
var server *http.Server

type userContextKey string

func Start() {
	createRoutes()

	server = &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}
	log.Println("Server starting on", os.Getenv("SERVER_ADDR"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
