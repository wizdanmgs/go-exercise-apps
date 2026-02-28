package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"file-upload-server/internal/service"
)

func TestUploadHandler(t *testing.T) {
	tempDir := t.TempDir()
	service := service.NewUploadService(tempDir)
	handler := NewUploadHandler(service)

	tests := []struct {
		name       string
		fileData   []byte
		expectCode int
	}{
		{
			name:       "valid upload",
			fileData:   append([]byte{0xFF, 0xD8, 0xFF}, make([]byte, 509)...),
			expectCode: http.StatusOK,
		},
		{
			name:       "invalid upload",
			fileData:   []byte("invalid"),
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, _ := writer.CreateFormFile("image", "test.jpg")
			if _, err := part.Write(tt.fileData); err != nil {
				t.Fatalf("failed to write file: %v", err.Error())
			}

			if err := writer.Close(); err != nil {
				t.Fatalf("failed to close writer: %v", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rr := httptest.NewRecorder()

			handler.Upload(rr, req)

			if rr.Code != tt.expectCode {
				t.Fatalf("expected status %d, got %d", tt.expectCode, rr.Code)
			}
		})
	}
}
