//go:generate npm run build

package main

import (
	"github.com/joho/godotenv"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/pkg/db"
)

func main() {
	godotenv.Load(".env")
	db.Connect()
	models.CreateTables()
	sessions.Start()
	handlers.Start()
}
