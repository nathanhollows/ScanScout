package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nathanhollows/Rapua/internal/filesystem"
	"github.com/nathanhollows/Rapua/internal/handlers/internal/middlewares"
)

func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	router = chi.NewRouter()
	router.Use(middleware.Compress(5))
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)

	router.Get("/", publicHomeHandler)

	// Session routes
	router.Get("/login", adminLoginHandler)
	router.Post("/login", adminLoginFormHandler)
	router.Get("/logout", adminLogoutHandler)
	router.Get("/register", adminRegisterHandler)
	router.Post("/register", adminRegisterFormHandler)

	// Scanning in routes
	router.Route("/s", func(r chi.Router) {
		r.Get("/{code:[A-z]{5}}", publicScanHandler)
		r.Post("/{code:[A-z]{5}}", publicScanPostHandler)
	})

	// Scanning out routes
	router.Route("/o", func(r chi.Router) {
		r.Get("/", publicScanOutHandler)
		r.Get("/{code:[A-z]{5}}", publicScanOutHandler)
		r.Post("/{code:[A-z]{5}}", publicScanOutPostHandler)
	})

	// Next location routes
	router.Get("/next", publicNextHandler)
	router.Post("/next", publicNextHandler)

	router.Route("/mylocations", func(r chi.Router) {
		r.Get("/", publicMyLocationsHandler)
		r.Get("/{code:[A-z]{5}}", publicSpecificLocationsHandler)
		r.Post("/{code:[A-z]{5}}", publicSpecificLocationsHandler)
	})

	router.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AdminAuthMiddleware)
		r.Use(middlewares.AdminCheckInstanceMiddleware)
		r.Get("/", adminDashboardHandler)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", adminLocationsHandler)
			r.Get("/new", adminLocationNewHandler)
			r.Post("/new", adminLocationNewPostHandler)
			r.Get("/{id}", adminLocationEditHandler)
			// r.Post("/{id}", adminLocationSaveHandler)
			// Disabled for now
			// r.Get("/qr/{id}.png", adminLocationQRHandler)
			r.Get("/qr-codes.zip", adminLocationQRZipHandler)
			r.Get("/posters.pdf", adminLocationPostersHandler)
		})

		r.Route("/teams", func(r chi.Router) {
			r.Get("/", adminTeamsHandler)
			r.Post("/add", adminTeamsAddHandler)
		})

		r.Route("/instances", func(r chi.Router) {
			r.Get("/", adminInstancesHandler)
			r.Post("/new", adminInstanceCreateHandler)
			r.Get("/{id}", adminInstancesHandler)
			r.Post("/{id}", adminInstancesHandler)
			r.Get("/{id}/switch", adminInstanceSwitchHandler)
			r.Post("/delete", adminInstanceDeleteHandler)
		})
	})

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "web/static"))}
	filesystem.FileServer(router, "/assets", filesDir)

	return router
}
