package services_test

import (
	"context"
	"os"
	"testing"
)

func TestCreateQRCodeImage(t *testing.T) {
	assetGen := NewAssetGenerator()

	tests := []struct {
		name      string
		path      string
		content   string
		options   []QRCodeOption
		expectErr bool
		cleanupFn func() // Function to clean up generated files if necessary.
	}{
		{
			name:      "Default options",
			path:      "default.png",
			content:   "Hello, QR!",
			expectErr: false,
			cleanupFn: func() { os.Remove("default.png") }, // Cleanup generated file after test
		},
		{
			name:    "Invalid format",
			path:    "invalid_format.qr",
			content: "Hello, QR!",
			options: []QRCodeOption{
				assetGen.WithFormat("bmp"),
			},
			expectErr: true,
		},
		{
			name:    "Valid PNG format",
			path:    "valid.png",
			content: "Hello, QR-1!",
			options: []QRCodeOption{
				assetGen.WithFormat("png"),
			},
			expectErr: false,
			cleanupFn: func() { os.Remove("valid.png") },
		},
		{
			name:    "Valid SVG format",
			path:    "valid.svg",
			content: "Hello, QR-2!",
			options: []QRCodeOption{
				assetGen.WithFormat("svg"),
			},
			expectErr: false,
			cleanupFn: func() { os.Remove("valid.svg") },
		},
		{
			name:    "Custom Colors",
			path:    "colored.svg",
			content: "Hello, QR-3!",
			options: []QRCodeOption{
				assetGen.WithFormat("svg"),
				assetGen.WithForeground("#FF0000"),
				assetGen.WithBackground("#00FF00"),
			},
			expectErr: false,
			cleanupFn: func() { os.Remove("colored.svg") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assetGen.CreateQRCodeImage(context.Background(), tt.path, tt.content, tt.options...)

			if (err != nil) != tt.expectErr {
				t.Fatalf("Expected error: %v, got: %v", tt.expectErr, err)
			}

			if tt.cleanupFn != nil {
				tt.cleanupFn()
			}
		})
	}
}
