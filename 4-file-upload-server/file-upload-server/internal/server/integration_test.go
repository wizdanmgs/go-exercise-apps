package server

import (
	"bytes"
	"encoding/json"
	"file-upload-server/internal/validator"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadAndServeIntegration(t *testing.T) {
	tempDir := t.TempDir()

	router := NewRouter(tempDir)

	// Step 1: Upload File
	imageData := validator.GenerateJPEG()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}

	_, err = part.Write(imageData)
	if err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	// Parse response
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	filename := resp["filename"]
	if filename == "" {
		t.Fatal("filename missing in response")
	}

	// Step 2: Verify File Saved
	fullPath := filepath.Join(tempDir, filename)
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("file not saved: %v", err)
	}

	// Step 3: Fetch Uploaded File
	getReq := httptest.NewRequest(http.MethodGet, "/uploads/"+filename, nil)
	getRR := httptest.NewRecorder()

	router.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected 200 when serving file, got %d", getRR.Code)
	}

	servedData, err := io.ReadAll(getRR.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(servedData, imageData) {
		t.Fatal("served file content does not match uploaded content")
	}
}
