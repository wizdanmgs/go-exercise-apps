package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"file-upload-server/internal/validator"
)

type UploadService struct {
	uploadDir string
}

func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{uploadDir: uploadDir}
}

func (s *UploadService) SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file type
	err := validator.ValidateImage(file)
	if err != nil {
		return "", err
	}

	// Reset file pointer after validation
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	destPath := filepath.Join(s.uploadDir, filename)
	if header.Size > 5<<20 { // 5MB max file
		return "", errors.New("file too large")
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := dst.Close(); err == nil {
			err = closeErr // only set err if no previous error
		}
	}()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filename, err
}
