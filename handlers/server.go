package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nathanhollows/ScanScout/filesystem"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/sessions"
)

var router *chi.Mux
var server *http.Server

func Start() {

	createRoutes()

	server = &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}
	fmt.Println(server.ListenAndServe())
}

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
		r.Get("/", adminActivityHandler)
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
		r.Get("/admin/instances", adminLocationsHandler)
	})

	workDir, _ := os.Getwd()
	filesDir := filesystem.Myfs{Dir: http.Dir(filepath.Join(workDir, "assets"))}
	filesystem.FileServer(router, "/assets", filesDir)

}

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session
		session, err := sessions.Get(r, "admin")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session.Values["user_id"] == nil {
			flash.Message{
				Title:   "Error",
				Message: "You must be logged in to access this page",
				Style:   flash.Error,
			}.Save(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func templateData(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"hxrequest": r.Header.Get("HX-Request") == "true",
		"layout":    "base",
	}
}

func render(w http.ResponseWriter, data map[string]interface{}, admin bool, patterns ...string) error {
	w.Header().Set("Content-Type", "text/html")

	baseDir := "templates/public/"
	if admin {
		baseDir = "templates/admin/"
	}

	err := parse(data, baseDir, patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), 0)
		log.Print("Template executing error: ", err)
	}
	return err
}

func parse(data map[string]interface{}, baseDir string, patterns ...string) *template.Template {
	// Format the title to include the app name
	if title, ok := data["title"].(string); ok {
		data["title"] = fmt.Sprintf("%s | %s", title, os.Getenv("APP_NAME"))
	}

	// Prepend the base directory to each pattern.
	for i, pattern := range patterns {
		patterns[i] = filepath.Join(baseDir, "pages", pattern+".html")
	}

	// Add the components dir to the patterns
	components, err := filepath.Glob(filepath.Join(baseDir, "components", "*.html"))
	if err != nil {
		log.Print("Error getting components: ", err)
	}
	patterns = append(patterns, components...)

	// Get the chosen layout
	if layout, ok := data["layout"].(string); ok {
		patterns = append(patterns, filepath.Join(baseDir, "layouts", layout+".html"))
	}

	// Create a new template, add any functions, and parse the files.
	return template.Must(template.New("base").Funcs(funcs).ParseFiles(patterns...))
}

var funcs = template.FuncMap{
	"html": func(v string) template.HTML {
		return template.HTML(v)
	},
	"upper": func(v string) string {
		return strings.ToUpper(v)
	},
	"lower": func(v string) string {
		return strings.ToLower(v)
	},
	"date": func(t time.Time) string {
		if t.Year() == time.Now().Year() {
			return t.Format("2 January")
		}
		return t.Format("2 January 2006")
	},
	"time": func(t time.Time) string {
		return t.Format("15:04")
	},
	"divide": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b)
	},
	"progress": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b) * 100
	},
	"add": func(a, b int) int {
		return a + b
	},
	"year": func() string {
		return time.Now().Format("2006")
	},
	"static": func(filename string) string {
		filename = strings.TrimPrefix(filename, "/")
		// get last modified time
		file, err := os.Stat("assets/" + filename)

		if err != nil {
			return "/assets/" + filename
		}

		modifiedtime := file.ModTime()
		return "/assets/" + filename + "?v=" + modifiedtime.Format("20060102150405")
	},
	// Convert a float to a duration and present it in a human readable format
	"toDuration": func(seconds float64) string {
		return time.Duration(int(seconds) * int(time.Second)).String()

	},
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
