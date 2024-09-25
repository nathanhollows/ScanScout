package services

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nathanhollows/Rapua/internal/helpers"
	go_qr "github.com/piglig/go-qr"
)

type QRCodeOptions struct {
	format     string
	scanType   string
	foreground string
	background string
}

type QRCodeOption func(*QRCodeOptions)

func (_ *assetGenerator) WithFormat(format string) QRCodeOption {
	return func(o *QRCodeOptions) {
		o.format = strings.ToLower(format)
	}
}

func (_ *assetGenerator) WithForeground(color string) QRCodeOption {
	return func(o *QRCodeOptions) {
		o.foreground = color
	}
}

func (_ *assetGenerator) WithBackground(color string) QRCodeOption {
	return func(o *QRCodeOptions) {
		o.background = color
	}
}

type AssetGenerator interface {
	// CreateQRCodeImage creates a QR code image with the given options
	// Supported options are:
	// - WithFormat(format string), where format is "png" or "svg"
	// - WithScanType(scanType string), where scanType is "in" or "out"
	// - WithForeground(color string), where color is a hex color code
	// - WithBackground(color string), where color is a hex color code
	CreateQRCodeImage(ctx context.Context, path string, content string, options ...QRCodeOption) (err error)
	// WithFormat sets the format of the QR code
	// Supported formats are "png" and "svg"
	WithFormat(format string) QRCodeOption
	// WithForeground sets the foreground color of the QR code
	WithForeground(color string) QRCodeOption
	// WithBackground sets the background color of the QR code
	WithBackground(color string) QRCodeOption

	// CreateArchive creates a zip archive from the given paths
	// Returns the path to the archive
	// Accepts a list of paths to files to add to the archive
	// Accepts an optional list of filenames to use for the files in the archive
	CreateArchive(ctx context.Context, paths []string) (path string, err error)
	CreatePDF(ctx context.Context, data []string) (string, error)
}

type assetGenerator struct{}

func NewAssetGenerator() AssetGenerator {
	return &assetGenerator{}
}

func (s *assetGenerator) CreateQRCodeImage(ctx context.Context, path string, content string, options ...QRCodeOption) (err error) {
	defaultOptions := &QRCodeOptions{
		format:     "png",
		scanType:   "in",
		foreground: "#000000",
		background: "#ffffff",
	}

	// Apply each option to the default options
	for _, o := range options {
		o(defaultOptions)
	}

	// Validate the options
	if defaultOptions.format != "png" && defaultOptions.format != "svg" {
		return errors.New(fmt.Sprintf("unsupported format: %s", defaultOptions.format))
	}

	qr, err := go_qr.EncodeText(content, go_qr.Medium)
	go_qr.MakeSegmentsOptimally(content, go_qr.Medium, 10, 27)
	config := go_qr.NewQrCodeImgConfig(20, 2)

	if defaultOptions.format == "png" {
		err := qr.PNG(config, path)
		if err != nil {
			return err
		}
	} else if defaultOptions.format == "svg" {
		err := qr.SVG(config, path, defaultOptions.background, defaultOptions.foreground)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *assetGenerator) CreateArchive(ctx context.Context, paths []string) (path string, err error) {
	// Create the file
	path = "assets/codes/" + helpers.NewCode(10) + "-" + fmt.Sprint(time.Now().UnixNano()) + ".zip"
	archive, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("could not create archive: %w", err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// Add each file to the zip
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return "", err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return "", err
		}

		header.Name = strings.TrimPrefix(path, "assets/codes/")
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func (s *assetGenerator) CreatePDF(ctx context.Context, data []string) (string, error) {
	return "", nil
}
