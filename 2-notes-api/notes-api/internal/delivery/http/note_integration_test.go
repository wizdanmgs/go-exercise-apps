package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"notes-api/internal/delivery/dto"
	delivery "notes-api/internal/delivery/http"
	"notes-api/internal/logger"
	"notes-api/internal/repository/memory"
	"notes-api/internal/usecase"
)

// setupTestServer builds full stack: repo -> usecase -> handler -> mux
func setupTestServer() *httptest.Server {
	logg := logger.New()
	repo := memory.NewMemoryRepository()
	uc := usecase.NewNoteUsecase(repo)
	handler := delivery.NewNoteHandler(uc, logg)

	r := chi.NewRouter()

	r.Route("/notes", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.GetAll)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetByID)
			r.Put("/", handler.Update)
			r.Delete("/", handler.Delete)
		})
	})

	return httptest.NewServer(r)
}

// ==== STANDARD MUX SERVER ====
// func setupTestServer() *httptest.Server {
// 	logg := logger.New()
// 	repo := memory.NewMemoryRepository()
// 	uc := usecase.NewNoteUsecase(repo)
// 	handler := h.NewNoteHandler(uc, logg)
//
// 	mux := http.NewServeMux()
//
// 	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodPost:
// 			handler.Create(w, r)
// 		case http.MethodGet:
// 			handler.GetAll(w, r)
// 		default:
// 			http.NotFound(w, r)
// 		}
// 	})
//
// 	mux.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodGet:
// 			handler.GetByID(w, r)
// 		case http.MethodPut:
// 			handler.Update(w, r)
// 		case http.MethodDelete:
// 			handler.Delete(w, r)
// 		default:
// 			http.NotFound(w, r)
// 		}
// 	})
//
// 	return httptest.NewServer(mux)
// }

func TestNotesIntegration(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := server.Client()

	// Create
	createReq := dto.CreateNoteRequest{
		ID:    "1",
		Title: "Integration Test",
	}
	body, _ := json.Marshal(createReq)

	resp, err := client.Post(server.URL+"/notes", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	// Get By ID
	resp, err = client.Get(server.URL + "/notes/1")
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var getResp dto.NoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&getResp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if getResp.Title != createReq.Title {
		t.Fatalf("expected %s, got %s", createReq.Title, getResp.Title)
	}

	// Update
	updated := dto.UpdateNoteRequest{
		Title: "Updated Title",
	}

	body, _ = json.Marshal(updated)

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/notes/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updatedResp dto.NoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedResp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if updatedResp.Title != updated.Title {
		t.Fatalf("expected %s, got %s", updated.Title, updatedResp.Title)
	}

	// Delete
	req, _ = http.NewRequest(http.MethodDelete, server.URL+"/notes/1", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Verify Deleted
	resp, err = client.Get(server.URL + "/notes/1")
	if err != nil {
		t.Fatalf("get after delete failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
