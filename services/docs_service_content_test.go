package services_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nathanhollows/Rapua/services"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Markdown to AST
func testDocs_MarkdownToAST(t *testing.T, markdown string) ast.Node {
	t.Helper()

	// Goldmark
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	// Parse markdown
	md := text.NewReader([]byte(markdown))
	var buf bytes.Buffer
	if err := gm.Convert([]byte(markdown), &buf); err != nil {
		t.Fatalf("failed to convert markdown: %v", err)
	}

	// Get AST
	node := gm.Parser().Parse(md)
	return node

}

// Make sure that all internal links are valid and point to an existing page.
func TestDocs_LinksResolve(t *testing.T) {
	dir := "../docs"
	docsService, err := services.NewDocsService(dir)
	if err != nil {
		t.Fatalf("failed to create DocsService: %v", err)
	}

	var walkPages func(pages []*services.DocPage)
	walkPages = func(pages []*services.DocPage) {
		for _, page := range pages {
			if len(page.Children) > 0 {
				walkPages(page.Children)
			}
			nodes := testDocs_MarkdownToAST(t, page.Content)
			ast.Walk(nodes, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
				if n.Kind() == ast.KindLink {
					link := n.(*ast.Link)
					dest := (string)(link.Destination)
					// Only check internal links
					if !strings.HasPrefix(dest, "/docs/") {
						return ast.WalkStop, nil
					}
					// Complain if the link doesn't resolve to a doc page
					if _, err := docsService.GetPage(dest); err != nil {
						t.Errorf("invalid link (%s) in /docs/%s", dest, page.Path)
					}
				}
				return ast.WalkContinue, nil
			})
		}
	}
	walkPages(docsService.Pages)
}
