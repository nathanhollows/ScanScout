package helpers

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var p = bluemonday.
	UGCPolicy().
	AddTargetBlankToFullyQualifiedLinks(true).
	// Allow iframe with class "enclave-object"
	AllowElementsMatching(regexp.MustCompile(`^iframe$`)).
	AllowAttrs("class").Matching(regexp.MustCompile(`\benclave-object\b`)).OnElements("iframe").
	AllowAttrs("src", "width", "height", "allow", "allowfullscreen", "frameborder").
	OnElements("iframe")

func SanitizeHTML(input []byte) []byte {
	return p.SanitizeBytes(input)
}
