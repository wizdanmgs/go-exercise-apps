package domain

// NoteRepository defines data persistence behavior.
// This belongs to domain because it defines business boundary.
type NoteRepository interface {
	Create(note Note) error
	GetAll() ([]Note, error)
	GetByID(id string) (Note, error)
	Update(id string, note Note) error
	Delete(id string) error
}
