//go:generate npm run build

package main

import (
	"github.com/joho/godotenv"
	"github.com/nathanhollows/ScanScout/handlers"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

func main() {
	godotenv.Load(".env")
	models.InitDB()
	sessions.Start()
	handlers.Start()
}
