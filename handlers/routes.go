package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nathanhollows/Rapua/filesystem"
)

func createRoutes() {
	router = chi.NewRouter()
	router.Use(middleware.Compress(5))
	router.Use(middleware.CleanPath)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)

	router.Get("/", publicHomeHandler)

	// Session routes
	router.Get("/login", adminLoginHandler)
	router.Post("/login", adminLoginPostHandler)
	router.Get("/logout", adminLogoutHandler)

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
		r.Post("/", adminLoginPostHandler)
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
		r.Use(adminAuthMiddleware)
		r.Get("/", adminDashboardHandler)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", adminLocationsHandler)
			r.Get("/new", adminLocationNewHandler)
			r.Post("/new", adminLocationSaveHandler)
			r.Get("/{id}", adminLocationEditHandler)
			r.Post("/{id}", adminLocationSaveHandler)
			// Disabled for now
			// r.Get("/qr/{id}.png", adminLocationQRHandler)
			r.Get("/qr/{id}.zip", adminLocationQRZipHandler)
			r.Get("/posters/{id}.pdf", adminLocationPostersHandler)
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
		})
	})

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "assets"))}
	filesystem.FileServer(router, "/assets", filesDir)
}
