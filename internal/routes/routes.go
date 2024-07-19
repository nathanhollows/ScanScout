package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nathanhollows/Rapua/internal/filesystem"
	"github.com/nathanhollows/Rapua/internal/handlers"
	admin "github.com/nathanhollows/Rapua/internal/handlers/admin"
	players "github.com/nathanhollows/Rapua/internal/handlers/players"
	public "github.com/nathanhollows/Rapua/internal/handlers/public"
	"github.com/nathanhollows/Rapua/internal/middlewares"
	"github.com/nathanhollows/Rapua/internal/services"
)

func SetupRouter(gameplayService *services.GameplayService, gameManagerService *services.GameManagerService) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Compress(5))
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)

	// Public routes
	setupPublicRoutes(router)

	// Player routes
	setupPlayerRoutes(router, gameplayService)

	// Admin routes
	setupAdminRoutes(router, gameManagerService)

	// Static files
	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "web/static"))}
	filesystem.FileServer(router, "/assets", filesDir)

	return router
}

// Setup the player routes
func setupPlayerRoutes(router chi.Router, gameplayService *services.GameplayService) {
	var playerHandler = players.NewPlayerHandler(gameplayService)

	// Home route
	// Takes a GET request to show the home page
	// Takes a POST request to submit the home page form
	router.Get("/", playerHandler.Home)
	router.Post("/", playerHandler.Home)

	// Show the next available locations
	router.Route("/next", func(r chi.Router) {
		r.Use(middlewares.TeamMiddleware)
		r.Get("/", playerHandler.Next)
		r.Post("/", playerHandler.Next)
	})

	// Check in to a location
	router.Route("/s", func(r chi.Router) {
		r.Use(middlewares.TeamMiddleware)
		r.Get("/{code:[A-z]{5}}", playerHandler.CheckIn)
		r.Post("/{code:[A-z]{5}}", playerHandler.ScanPost)
	})

	// Check out of a location
	router.Route("/o", func(r chi.Router) {
		r.Use(middlewares.TeamMiddleware)
		r.Get("/", playerHandler.ScanOut)
		r.Get("/{code:[A-z]{5}}", playerHandler.ScanOut)
		r.Post("/{code:[A-z]{5}}", playerHandler.ScanOutPost)
	})
}

func setupPublicRoutes(router chi.Router) {

	publicHandler := public.NewPublicHandler()

	router.Get("/home", publicHandler.Index)

	router.Route("/login", func(r chi.Router) {
		r.Get("/", publicHandler.Login)
		r.Post("/", publicHandler.LoginPost)
	})
	router.Get("/logout", publicHandler.Logout)
	router.Route("/register", func(r chi.Router) {
		r.Get("/", publicHandler.Register)
		r.Post("/", publicHandler.RegisterPost)
	})

	router.Route("/mylocations", func(r chi.Router) {
		r.Use(middlewares.TeamMiddleware)
		r.Get("/", handlers.PublicMyLocationsHandler)
		r.Get("/{code:[A-z]{5}}", handlers.PublicSpecificLocationsHandler)
		r.Post("/{code:[A-z]{5}}", handlers.PublicSpecificLocationsHandler)
	})
}

func setupAdminRoutes(router chi.Router, gameManagerService *services.GameManagerService) {
	var adminHandler = admin.NewAdminHandler(gameManagerService)

	router.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AdminAuthMiddleware)
		r.Use(middlewares.AdminCheckInstanceMiddleware)

		r.Get("/", adminHandler.Activity)

		r.Route("/locations", func(r chi.Router) {
			r.Get("/", adminHandler.Locations)
			r.Get("/new", adminHandler.LocationNew)
			r.Post("/new", adminHandler.LocationNewPost)
			r.Get("/{id}", handlers.AdminLocationEditHandler)
			r.Post("/{id}", handlers.AdminLocationEditPostHandler)
			// Disabled for now
			// r.Get("/qr/{id}.png", handlers.AdminLocationQRHandler)
			r.Get("/qr-codes.zip", handlers.AdminLocationQRZipHandler)
			r.Get("/posters.pdf", handlers.AdminLocationPostersHandler)
		})

		r.Route("/teams", func(r chi.Router) {
			r.Get("/", adminHandler.Teams)
			r.Post("/add", adminHandler.TeamsAdd)
		})

		r.Route("/instances", func(r chi.Router) {
			r.Get("/", adminHandler.Instances)
			r.Post("/new", adminHandler.InstancesCreate)
			r.Get("/{id}", adminHandler.Instances)
			r.Post("/{id}", adminHandler.Instances)
			r.Get("/{id}/switch", adminHandler.InstanceSwitch)
			r.Post("/delete", adminHandler.InstanceDelete)
		})
	})
}
