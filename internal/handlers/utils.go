package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type TemplateDir string

const (
	AdminDir  TemplateDir = "admin"
	PlayerDir TemplateDir = "players"
	PublicDir TemplateDir = "public"
)

func SetDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func TemplateData(r *http.Request) map[string]interface{} {
	user, ok := r.Context().Value(contextkeys.UserIDKey).(*models.User)
	data := map[string]interface{}{
		"hxrequest": r.Header.Get("HX-Request") == "true",
		"layout":    "base",
	}
	if ok && user != nil {
		data["user"] = user
		data["instances"] = user.Instances
	}
	return data
}

func Render(
	w http.ResponseWriter,
	data map[string]interface{},
	templateDir TemplateDir,
	patterns ...string,
) error {
	w.Header().Set("Content-Type", "text/html")

	baseDir := "web/templates/" + string(templateDir) + "/"

	err := parse(data, baseDir, patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		slog.Error("executing template", "err", err)
	}
	return err
}

func RenderHTMX(
	w http.ResponseWriter,
	data map[string]interface{},
	templateDir TemplateDir,
	patterns ...string,
) error {
	w.Header().Set("Content-Type", "text/html")

	baseDir := "web/templates/" + string(templateDir) + "/"
	tmpl := parse(data, baseDir, patterns...)
	err := tmpl.ExecuteTemplate(w, "content", data)
	if err != nil {
		slog.Error("executing template", "err", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
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
	"getEnv": func(key string) string {
		return os.Getenv(key)
	},
	// Render a string as HTML
	"html": func(v string) template.HTML {
		return template.HTML(v)
	},
	// Convert a string to uppercase
	"upper": func(v string) string {
		return strings.ToUpper(v)
	},
	// Convert a string to lowercase
	"lower": func(v string) string {
		return strings.ToLower(v)
	},
	// Format a time.Time to a human readable date
	"date": func(t time.Time) string {
		if t.Year() == time.Now().Year() {
			return t.Format("2 January")
		}
		return t.Format("2 January 2006")
	},
	// Format a time.Time to a human readable time
	"time": func(t time.Time) string {
		return t.Format("15:04")
	},
	// Easy division
	"divide": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b)
	},
	// Replace newlines with <br> tags
	"nl2br": func(s string) template.HTML {
		return template.HTML(strings.Replace(s, "\n", "<br>", -1))
	},
	// Calculate the progress percentage
	"progress": func(a, b int) float32 {
		if a == 0 || b == 0 {
			return 0
		}
		return float32(a) / float32(b) * 100
	},
	// Adds two integers
	"add": func(a, b int) int {
		return a + b
	},
	// Return the current year
	"year": func() string {
		return time.Now().Format("2006")
	},
	// Link to a static file with a cache busting query string
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
	// Convert markdown (string) to HTML
	"md": func(s string) template.HTML {
		md := goldmark.New(
			goldmark.WithExtensions(
				extension.Strikethrough,
				extension.Linkify,
				extension.TaskList,
				extension.Typographer,
				enclave.New(
					&enclave.Config{},
				),
			),
			goldmark.WithParserOptions(),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
			),
		)

		var buf bytes.Buffer
		if err := md.Convert([]byte(s), &buf); err != nil {
			slog.Error("converting markdown to HTML", "err", err)
			return template.HTML("Error rendering markdown to HTML")
		}

		return template.HTML(helpers.SanitizeHTML(buf.Bytes()))
	},
	// Sequence returns a slice of integers
	// Accepts 1, 2, or 3 int arguments
	"sequence": func(args ...int) []int {
		switch len(args) {
		case 1:
			return make([]int, args[0])
		case 2:
			return make([]int, args[1]-args[0])
		case 3:
			s := make([]int, args[2]-args[0])
			for i := range s {
				s[i] = args[0] + i
			}
		}
		return []int{}
	},
}
