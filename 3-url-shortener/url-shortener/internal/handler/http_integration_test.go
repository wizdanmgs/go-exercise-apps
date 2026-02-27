package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-shortener/internal/handler"
	"url-shortener/internal/service"
	"url-shortener/internal/store"

	"github.com/go-chi/chi/v5"
)

func setupTestServer() *httptest.Server {
	store := store.NewMemoryStore()
	shortener := service.NewShortener(store)
	h := handler.NewHandler(shortener)

	r := chi.NewRouter()
	r.Post("/shorten", h.Create)
	r.Get("/{code}", h.Redirect)

	return httptest.NewServer(r)
}

func TestIntegration_CreateAndRedirect(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// Step 1: Create short URL
	reqBody := map[string]interface{}{
		"url": "https://example.com",
		"ttl": 3600,
	}

	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(ts.URL+"/shorten", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to call shorten: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("failed to close body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var createResp struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createResp.Code == "" {
		t.Fatalf("expected non-empty code")
	}

	// Step 2: Follow redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	redirectResp, err := client.Get(ts.URL + "/" + createResp.Code)
	if err != nil {
		t.Fatalf("expected 302, got %d", redirectResp.StatusCode)
	}
	defer func() {
		if err := redirectResp.Body.Close(); err != nil {
			t.Fatalf("failed to close body: %v", err)
		}
	}()

	location := redirectResp.Header.Get("Location")
	if location != "https://example.com" {
		t.Fatalf("expected redirect location: %s", location)
	}
}

func TestIntegration_ExpiredLink(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	reqBody := map[string]any{
		"url": "https://expired.com",
		"ttl": -1, // already expired
	}

	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(ts.URL+"/shorten", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to call shorten: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("failed to close body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var createResp struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	redirectResp, err := client.Get(ts.URL + "/" + createResp.Code)
	if err != nil {
		t.Fatalf("failed to call redirect: %v", err)
	}
	defer func() {
		if err := redirectResp.Body.Close(); err != nil {
			t.Fatalf("failed to close body: %v", err)
		}
	}()

	if redirectResp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for expired link, got %d", redirectResp.StatusCode)
	}
}
