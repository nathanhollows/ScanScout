package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gomarkdown/markdown"
	"github.com/nathanhollows/Rapua/models"
)

func templateData(r *http.Request) map[string]interface{} {
	user, ok := r.Context().Value(models.UserIDKey).(*models.User)
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

func render(
	w http.ResponseWriter,
	data map[string]interface{},
	admin bool,
	patterns ...string,
) error {
	w.Header().Set("Content-Type", "text/html")

	baseDir := "templates/public/"
	if admin {
		baseDir = "templates/admin/"
	}

	err := parse(data, baseDir, patterns...).ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
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
		// Convert markdown to HTML
		content := []byte(s)
		return template.HTML(markdown.ToHTML(content, nil, nil))
	},
}
