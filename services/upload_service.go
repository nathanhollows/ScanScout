package services

import (
	"context"
	"mime/multipart"
)

type FileUploadService interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, instanceID string) (UploadResult, error)
	DeleteFile(ctx context.Context, fileName string, instanceID string) error
}

type UploadResult struct {
	URL      string
	FileName string
	FileSize int64
}
