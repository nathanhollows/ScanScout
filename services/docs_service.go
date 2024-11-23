package services

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// DocPage represents a single documentation page.
type DocPage struct {
	Title    string
	Order    int
	Path     string
	Content  string
	URL      string
	Headings []Heading
	Children []*DocPage
}

// Heading represents a section heading within a doc page.
type Heading struct {
	Level int
	Text  string
	ID    string
}

// DocsService handles loading and providing documentation content.
type DocsService struct {
	DocsDir string
	Pages   []*DocPage
}

// NewDocsService creates a new instance of DocsService.
func NewDocsService(docsDir string) (*DocsService, error) {
	service := &DocsService{
		DocsDir: docsDir,
	}

	if err := service.loadDocs(); err != nil {
		return nil, err
	}

	return service, nil
}

// loadDocs loads and parses all Markdown files in the DocsDir.
func (ds *DocsService) loadDocs() error {
	var pages []*DocPage

	err := filepath.Walk(ds.DocsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Markdown files
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Get the relative path for URL generation
		relativePath, err := filepath.Rel(ds.DocsDir, path)
		if err != nil {
			return err
		}
		relativePath = filepath.ToSlash(relativePath) // Ensure consistent path separators

		// Read file content
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Split front matter and content
		parts := strings.SplitN(string(data), "---", 3)
		if len(parts) < 3 {
			return nil // Skip files without proper front matter
		}

		// Parse YAML front matter
		var meta struct {
			Title   string `yaml:"title"`
			Order   int    `yaml:"order"`
			Sidebar *bool  `yaml:"sidebar"` // New field to capture sidebar visibility
		}
		if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
			return err
		}

		// Extract headings for ToC
		headings := extractHeadings(parts[2])

		// Create DocPage
		page := &DocPage{
			Title:    meta.Title,
			Order:    meta.Order,
			Path:     relativePath,
			URL:      "/docs/" + strings.TrimSuffix(relativePath, ".md"),
			Content:  parts[2],
			Headings: headings,
		}

		pages = append(pages, page)
		return nil
	})

	if err != nil {
		return err
	}

	// Build the page hierarchy
	ds.Pages = buildHierarchy(pages)
	return nil
}

// extractHeadings extracts headings from Markdown content.
func extractHeadings(content string) []Heading {
	var headings []Heading
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			level := len(strings.SplitN(line, " ", 2)[0])
			text := strings.TrimSpace(strings.TrimLeft(line, "# "))
			id := strings.ReplaceAll(strings.ToLower(text), " ", "-")
			headings = append(headings, Heading{
				Level: level,
				Text:  text,
				ID:    id,
			})
		}
	}
	return headings
}

// buildHierarchy organizes pages into a tree based on their paths.
func buildHierarchy(pages []*DocPage) []*DocPage {
	root := make(map[string]*DocPage)

	for _, page := range pages {
		parts := strings.Split(page.Path, "/")
		addToTree(root, parts, page, 0)
	}

	// Convert map to slice and sort
	var rootPages []*DocPage
	for _, page := range root {
		rootPages = append(rootPages, page)
	}

	sortPages(rootPages)
	return rootPages
}

func addToTree(node map[string]*DocPage, parts []string, page *DocPage, depth int) {
	if depth >= len(parts) {
		return
	}
	key := parts[depth]
	if existing, ok := node[key]; ok {
		// Existing node, proceed to next depth
		if depth == len(parts)-1 {
			// Leaf node
			existing.Title = page.Title
			existing.Order = page.Order
			existing.Content = page.Content
			existing.Headings = page.Headings
			existing.URL = page.URL
		} else {
			// Check if page is index.md for this directory
			if filepath.Base(page.Path) == "index.md" && depth == len(parts)-2 {
				// Update directory node with index.md's Title, Order, Content, etc.
				existing.Title = page.Title
				existing.Order = page.Order
				existing.Content = page.Content
				existing.Headings = page.Headings
				existing.URL = page.URL
			}
			if existing.Children == nil {
				existing.Children = []*DocPage{}
			}
			childMap := make(map[string]*DocPage)
			for _, child := range existing.Children {
				childMap[filepath.Base(child.Path)] = child
			}
			addToTree(childMap, parts, page, depth+1)
			existing.Children = mapToSlice(childMap)
			sortPages(existing.Children)
		}
	} else {
		// Create a new node
		newPage := &DocPage{
			Title: key,
			Path:  strings.Join(parts[:depth+1], "/"),
			URL:   "/docs/" + strings.Join(parts[:depth+1], "/"),
			Order: 9999, // Default order for directories without specified order
		}
		if depth == len(parts)-1 {
			// Leaf node
			newPage.Title = page.Title
			newPage.Order = page.Order
			newPage.Content = page.Content
			newPage.Headings = page.Headings
			newPage.URL = page.URL
		} else {
			// Check if page is index.md for this directory
			if filepath.Base(page.Path) == "index.md" && depth == len(parts)-2 {
				// Update directory node with index.md's Title, Order, Content, etc.
				newPage.Title = page.Title
				newPage.Order = page.Order
				newPage.Content = page.Content
				newPage.Headings = page.Headings
				newPage.URL = page.URL
			}
			childMap := make(map[string]*DocPage)
			addToTree(childMap, parts, page, depth+1)
			newPage.Children = mapToSlice(childMap)
			sortPages(newPage.Children)
		}
		node[key] = newPage
	}
}

func mapToSlice(m map[string]*DocPage) []*DocPage {
	var slice []*DocPage
	for _, v := range m {
		slice = append(slice, v)
	}
	return slice
}

func sortPages(pages []*DocPage) {
	sort.SliceStable(pages, func(i, j int) bool {
		if strings.Contains(pages[i].Path, "index.md") {
			return true
		}
		if strings.Contains(pages[j].Path, "index.md") {
			return false
		}
		return pages[i].Order < pages[j].Order
	})
	for _, page := range pages {
		if len(page.Children) > 0 {
			sortPages(page.Children)
		}
	}
}

func (ds *DocsService) GetPage(urlPath string) (*DocPage, error) {
	// Normalize the path
	trimmedPath := strings.TrimPrefix(urlPath, "/docs")
	trimmedPath = strings.TrimSuffix(trimmedPath, "/")

	if trimmedPath == "" {
		trimmedPath = "/index"
	}

	parts := strings.Split(trimmedPath, "/")[1:] // Skip the empty string at index 0

	var currentPages []*DocPage = ds.Pages
	var foundPage *DocPage

	for i, part := range parts {
		found := false
		for _, page := range currentPages {
			pageBaseName := strings.TrimSuffix(filepath.Base(page.Path), ".md")
			if pageBaseName == part {
				found = true
				foundPage = page
				currentPages = page.Children
				break
			}
		}
		if !found {
			// At the last segment, attempt to find index.md or first ordered page
			if i == len(parts)-1 {
				// Try to find index.md
				for _, page := range currentPages {
					pageBaseName := strings.TrimSuffix(filepath.Base(page.Path), ".md")
					if pageBaseName == "index" {
						return page, nil
					}
				}
				// If no index.md, find the page with lowest order
				if len(currentPages) > 0 {
					sortPages(currentPages)
					return currentPages[0], nil // Return the page with lowest order
				}
			}
			return nil, os.ErrNotExist
		}
	}
	if foundPage != nil {
		return foundPage, nil
	}
	return nil, os.ErrNotExist
}
