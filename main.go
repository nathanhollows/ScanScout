//go:generate npm run build

package main

import (
	"github.com/joho/godotenv"
	"github.com/nathanhollows/Rapua/handlers"
	"github.com/nathanhollows/Rapua/models"
	"github.com/nathanhollows/Rapua/sessions"
)

func main() {
	godotenv.Load(".env")
	models.InitDB()
	sessions.Start()
	handlers.Start()
}
