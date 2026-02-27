package memory

import (
	"sync"

	"notes-api/internal/domain"
)

// MemoryRepository is an in-memory implementation
// of the domain.NoteRepository interface.
type MemoryRepository struct {
	mu    sync.RWMutex
	notes map[string]domain.Note
}

// NewMemoryRepository initializes storage.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		notes: make(map[string]domain.Note),
	}
}

func (r *MemoryRepository) Create(note domain.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.notes[note.ID] = note
	return nil
}

func (r *MemoryRepository) GetAll() ([]domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Note
	for _, n := range r.notes {
		result = append(result, n)
	}
	return result, nil
}

func (r *MemoryRepository) GetByID(id string) (domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, ok := r.notes[id]
	if !ok {
		return domain.Note{}, domain.ErrNotFound
	}
	return note, nil
}

func (r *MemoryRepository) Update(id string, note domain.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return domain.ErrNotFound
	}

	note.ID = id
	r.notes[id] = note
	return nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return domain.ErrNotFound
	}

	delete(r.notes, id)
	return nil
}
