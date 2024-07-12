//go:generate npm run build

package main

import (
	"github.com/joho/godotenv"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/server"
	"github.com/nathanhollows/Rapua/internal/sessions"
	"github.com/nathanhollows/Rapua/pkg/db"
)

func main() {
	godotenv.Load(".env")
	db.Connect()
	models.CreateTables()
	sessions.Start()
	server.Start()
}
