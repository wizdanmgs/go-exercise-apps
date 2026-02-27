package memory

import (
	"errors"
	"testing"

	"notes-api/internal/domain"
)

func TestMemoryRepository_CRUD(t *testing.T) {
	repo := NewMemoryRepository()

	note := domain.Note{
		ID:    "1",
		Title: "Test",
	}

	// Create
	if err := repo.Create(note); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// GetByID
	result, err := repo.GetByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "1" {
		t.Fatalf("expected ID 1, got %s", result)
	}

	// Update
	note.Title = "Updated"
	if err := repo.Update("1", note); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := repo.GetByID("1")
	if updated.Title != "Updated" {
		t.Fatalf("update failed")
	}

	// Delete
	if err := repo.Delete("1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.GetByID("1")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound")
	}
}
