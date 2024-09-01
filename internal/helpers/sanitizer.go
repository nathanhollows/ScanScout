package helpers

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var p = bluemonday.
	UGCPolicy().
	AddTargetBlankToFullyQualifiedLinks(true).
	// Allow iframe with class "enclave-object"
	AllowAttrs("class").Matching(regexp.MustCompile(`\benclave-object\b`)).OnElements("iframe").
	AllowAttrs("src", "width", "height", "allow", "allowfullscreen", "frameborder").
	OnElements("iframe").
	// Allow input with type "checkbox", remove disabled attribute
	AllowAttrs("type").Matching(regexp.MustCompile(`\bcheckbox\b`)).OnElements("input")

func SanitizeHTML(input []byte) []byte {
	return input
}
