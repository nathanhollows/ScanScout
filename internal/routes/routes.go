package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nathanhollows/Rapua/internal/filesystem"
	"github.com/nathanhollows/Rapua/internal/handlers"
	"github.com/nathanhollows/Rapua/internal/middlewares"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	router = chi.NewRouter()
	router.Use(middleware.Compress(5))
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)

	router.Get("/", handlers.PublicHomeHandler)

	// Session routes
	router.Get("/login", handlers.AdminLoginHandler)
	router.Post("/login", handlers.AdminLoginFormHandler)
	router.Get("/logout", handlers.AdminLogoutHandler)
	router.Get("/register", handlers.AdminRegisterHandler)
	router.Post("/register", handlers.AdminRegisterFormHandler)

	// Scanning in routes
	router.Route("/s", func(r chi.Router) {
		r.Get("/{code:[A-z]{5}}", handlers.PublicScanHandler)
		r.Post("/{code:[A-z]{5}}", handlers.PublicScanPostHandler)
	})

	// Scanning out routes
	router.Route("/o", func(r chi.Router) {
		r.Get("/", handlers.PublicScanOutHandler)
		r.Get("/{code:[A-z]{5}}", handlers.PublicScanOutHandler)
		r.Post("/{code:[A-z]{5}}", handlers.PublicScanOutPostHandler)
	})

	// Next location routes
	router.Get("/next", handlers.PublicNextHandler)
	router.Post("/next", handlers.PublicNextHandler)

	router.Route("/mylocations", func(r chi.Router) {
		r.Get("/", handlers.PublicMyLocationsHandler)
		r.Get("/{code:[A-z]{5}}", handlers.PublicSpecificLocationsHandler)
		r.Post("/{code:[A-z]{5}}", handlers.PublicSpecificLocationsHandler)
	})

	router.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AdminAuthMiddleware)
		r.Use(middlewares.AdminCheckInstanceMiddleware)
		r.Get("/", handlers.AdminDashboardHandler)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", handlers.AdminLocationsHandler)
			r.Get("/new", handlers.AdminLocationNewHandler)
			r.Post("/new", handlers.AdminLocationNewPostHandler)
			r.Get("/{id}", handlers.AdminLocationEditHandler)
			// r.Post("/{id}", handlers.AdminLocationSaveHandler)
			// Disabled for now
			// r.Get("/qr/{id}.png", handlers.AdminLocationQRHandler)
			r.Get("/qr-codes.zip", handlers.AdminLocationQRZipHandler)
			r.Get("/posters.pdf", handlers.AdminLocationPostersHandler)
		})

		r.Route("/teams", func(r chi.Router) {
			r.Get("/", handlers.AdminTeamsHandler)
			r.Post("/add", handlers.AdminTeamsAddHandler)
		})

		r.Route("/instances", func(r chi.Router) {
			r.Get("/", handlers.AdminInstancesHandler)
			r.Post("/new", handlers.AdminInstanceCreateHandler)
			r.Get("/{id}", handlers.AdminInstancesHandler)
			r.Post("/{id}", handlers.AdminInstancesHandler)
			r.Get("/{id}/switch", handlers.AdminInstanceSwitchHandler)
			r.Post("/delete", handlers.AdminInstanceDeleteHandler)
		})
	})

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "web/static"))}
	filesystem.FileServer(router, "/assets", filesDir)

	return router
}
