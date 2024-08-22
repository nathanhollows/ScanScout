//go:generate npm run build

package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/server"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/pkg/db"
)

func main() {
	godotenv.Load(".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db.MustOpen()
	models.CreateTables(logger)
	sessions.Start()

	server.Start(logger)
}
