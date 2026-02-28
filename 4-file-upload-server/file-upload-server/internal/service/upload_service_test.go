package service

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func TestSaveFile(t *testing.T) {
	tempDir := t.TempDir()

	service := NewUploadService(tempDir)

	validJPEG := append([]byte{0xFF, 0xD8, 0xFF}, make([]byte, 509)...)
	invalidFile := []byte("not an image")

	tests := []struct {
		name      string
		content   []byte
		expectErr bool
	}{
		{
			name:      "valid image",
			content:   validJPEG,
			expectErr: false,
		},
		{
			name:      "invalid image",
			content:   invalidFile,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			mockFile := &mockMultipartFile{
				Reader: bytes.NewReader(tt.content),
			}

			header := &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     int64(len(tt.content)),
			}

			filename, err := service.SaveFile(
				mockFile, header,
			)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			fullPath := filepath.Join(tempDir, filename)

			if _, err := os.Stat(fullPath); err != nil {
				t.Fatalf("file not saved: %v", err)
			}
		})
	}
}
