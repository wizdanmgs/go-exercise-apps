package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"notes-api/internal/delivery/dto"
	"notes-api/internal/domain"
	"notes-api/internal/logger"
	"notes-api/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// NoteHandler handles HTTP requests.
// It depends on usecase, not repository.
type NoteHandler struct {
	usecase *usecase.NoteUsecase
	logger  *logger.Logger
}

// NewNoteHandler injects usecase dependency.
func NewNoteHandler(u *usecase.NoteUsecase, log *logger.Logger) *NoteHandler {
	return &NoteHandler{usecase: u, logger: log}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// Create handles POST /notes
func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid_request_body", "error", err)
		respondJSON(w, mapErrorToStatus(err), map[string]string{"error": "invalid body"})
		return
	}

	note := req.ToDomain()

	if err := h.usecase.Create(note); err != nil {
		h.logger.Error("failed_create_note", "error", err)
		respondJSON(w, mapErrorToStatus(err), map[string]string{"error": err.Error()})
		return
	}

	resp := dto.ToResponse(note)
	h.logger.Info("note_created", "note_id", resp.ID)
	respondJSON(w, http.StatusCreated, resp)
}

// GetAll handles GET /notes
func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.usecase.GetAll()
	if err != nil {
		h.logger.Error("failed_get_notes", "error", err)
		respondJSON(w, mapErrorToStatus(err), nil)
		return
	}

	var responses []dto.NoteResponse
	for _, n := range notes {
		responses = append(responses, dto.ToResponse(n))
	}

	h.logger.Info("notes_fetched")
	respondJSON(w, http.StatusOK, responses)
}

// GetByID handles GET /notes/{id}
func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// id := strings.TrimPrefix(r.URL.Path, "/notes/")

	note, err := h.usecase.GetByID(id)
	if err != nil {
		h.logger.Error("failed_get_note", "error", err)
		respondJSON(w, mapErrorToStatus(err), map[string]string{"error": err.Error()})
		return
	}

	resp := dto.ToResponse(note)
	h.logger.Info("note_fetched", "note_id", resp.ID)
	respondJSON(w, http.StatusOK, resp)
}

// Update handles PUT /notes/{id}
func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// id := strings.TrimPrefix(r.URL.Path, "/notes/")

	var req dto.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid_request_body", "error", err)
		respondJSON(w, mapErrorToStatus(err), map[string]string{
			"error": "invalid body",
		})
		return
	}

	note := domain.Note{
		ID:    id,
		Title: req.Title,
	}

	if err := h.usecase.Update(id, note); err != nil {
		h.logger.Error("failed_update_note", "error", err)
		respondJSON(w, mapErrorToStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	resp := dto.ToResponse(note)
	h.logger.Info("note_updated", "note_id", resp.ID)
	respondJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /notes/{id}
func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// id := strings.TrimPrefix(r.URL.Path, "/notes/")

	if err := h.usecase.Delete(id); err != nil {
		h.logger.Error("failed_delete_note", "error", err)
		respondJSON(w, mapErrorToStatus(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("note_deleted", "note_id", id)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "deleted",
	})
}
