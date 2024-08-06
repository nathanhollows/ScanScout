package helpers

import "testing"

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "Valid HTML",
			input: []byte("<p>Hello, <a href='https://example.com'>World</a></p>"),
			want:  []byte("<p>Hello, <a href=\"https://example.com\" rel=\"nofollow noopener\" target=\"_blank\">World</a></p>"),
		},
		{
			name:  "Invalid HTML",
			input: []byte("<script>alert('Hello, World!')</script>"),
			want:  []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeHTML(tt.input)
			if string(got) != string(tt.want) {
				t.Errorf("SanitizeHTML() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
