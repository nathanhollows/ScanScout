//go:generate npm run build

package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/nathanhollows/Rapua/db"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/server"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

func main() {
	godotenv.Load(".env")
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	initialiseFolders(logger)
	db.MustOpen()
	models.CreateTables(logger)
	sessions.Start()

	server.Start(logger)
}

func initialiseFolders(logger *slog.Logger) {
	folders := []string{
		"assets/",
		"assets/codes/",
		"assets/codes/png/",
		"assets/codes/svg/",
		"assets/fonts/",
		"assets/posters/"}

	for _, folder := range folders {
		_, err := os.Stat(folder)
		if err != nil {
			// Attempt to create the directory
			err = os.MkdirAll(folder, 0755)
			if err != nil {
				logger.Error("could not create directory", "folder", folder, "error", err)
				panic("Directory '" + folder + "' does not exist and could not be created")
			}
		}
	}
}
