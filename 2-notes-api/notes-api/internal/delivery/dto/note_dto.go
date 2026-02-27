package dto

import "notes-api/internal/domain"

// CreateNoteRequest represents incoming create request body
type CreateNoteRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// UpdateNoteRequest represents update request body.
type UpdateNoteRequest struct {
	Title string `json:"title"`
}

// NoteResponse represents outgoing response body.
type NoteResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ToDomain converts CreateNoteRequest to domain model.
func (r CreateNoteRequest) ToDomain() domain.Note {
	return domain.Note{
		ID:    r.ID,
		Title: r.Title,
	}
}

// ToResponse converts domain model to response DTO.
func ToResponse(n domain.Note) NoteResponse {
	return NoteResponse{
		ID:    n.ID,
		Title: n.Title,
	}
}
