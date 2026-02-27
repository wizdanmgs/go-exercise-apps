package usecase

import (
	"notes-api/internal/domain"
)

// NoteUsecase contains business logic.
// It depends only on domain interfaces.
type NoteUsecase struct {
	repo domain.NoteRepository
}

// NewNoteUsecase injects repository dependency.
func NewNoteUsecase(repo domain.NoteRepository) *NoteUsecase {
	return &NoteUsecase{repo: repo}
}

// Create validates and creates a note.
func (u *NoteUsecase) Create(note domain.Note) error {
	if note.ID == "" || note.Title == "" {
		return domain.ErrInvalidInput
	}

	return u.repo.Create(note)
}

// GetAll retrieves all notes.
func (u *NoteUsecase) GetAll() ([]domain.Note, error) {
	return u.repo.GetAll()
}

// GetByID retrieves a note by ID.
func (u *NoteUsecase) GetByID(id string) (domain.Note, error) {
	return u.repo.GetByID(id)
}

// Update updates a note.
func (u *NoteUsecase) Update(id string, note domain.Note) error {
	if note.Title == "" {
		return domain.ErrInvalidInput
	}
	return u.repo.Update(id, note)
}

// Delete removes a note.
func (u *NoteUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}
