package helpers

import "testing"

func TestNewCode(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "Short Code",
			length: 6,
		},
		{
			name:   "Long Code",
			length: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewCode(tt.length)
			if len(code) != tt.length {
				t.Errorf("NewCode() = %v, want %v", len(code), tt.length)
			}
		})
	}
}
