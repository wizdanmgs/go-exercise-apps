package usecase

import (
	"errors"
	"testing"

	"notes-api/internal/domain"
)

// mockRepo implements domain.NoteRepository for testing.
type mockRepo struct {
	createFn  func(note domain.Note) error
	getAllFn  func() ([]domain.Note, error)
	getByIDFn func(id string) (domain.Note, error)
	updateFn  func(id string, note domain.Note) error
	deleteFn  func(id string) error
}

func (m *mockRepo) Create(note domain.Note) error {
	return m.createFn(note)
}

func (m *mockRepo) GetAll() ([]domain.Note, error) {
	return m.getAllFn()
}

func (m *mockRepo) GetByID(id string) (domain.Note, error) {
	return m.getByIDFn(id)
}

func (m *mockRepo) Update(id string, note domain.Note) error {
	return m.updateFn(id, note)
}

func (m *mockRepo) Delete(id string) error {
	return m.deleteFn(id)
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		note    domain.Note
		repoErr error
		wantErr error
	}{
		{
			name: "success",
			note: domain.Note{
				ID:    "1",
				Title: "Test",
			},
			wantErr: nil,
		},
		{
			name: "missing id",
			note: domain.Note{
				Title: "Test",
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "repo failure",
			note: domain.Note{
				ID:    "1",
				Title: "Test",
			},
			repoErr: domain.ErrDb,
			wantErr: domain.ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := &mockRepo{
				createFn: func(note domain.Note) error {
					return tt.repoErr
				},
			}

			uc := NewNoteUsecase(mock)

			err := uc.Create(tt.note)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name    string
		note    []domain.Note
		repoErr error
		wantErr error
	}{
		{
			name: "success",
			note: []domain.Note{
				{
					ID:    "1",
					Title: "Test 1",
				},
				{
					ID:    "2",
					Title: "Test 2",
				},
			},
			wantErr: nil,
		},
		{
			name: "repo failure",
			note: []domain.Note{
				{
					ID:    "1",
					Title: "Test 1",
				},
				{
					ID:    "2",
					Title: "Test 2",
				},
			},
			repoErr: domain.ErrDb,
			wantErr: domain.ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := &mockRepo{
				getAllFn: func() ([]domain.Note, error) {
					return tt.note, tt.repoErr
				},
			}

			uc := NewNoteUsecase(mock)

			_, err := uc.GetAll()

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repoErr error
		wantErr error
	}{
		{
			name: "success",
			id:   "1",
		},
		{
			name:    "not found",
			id:      "1",
			repoErr: domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := &mockRepo{
				getByIDFn: func(id string) (domain.Note, error) {
					return domain.Note{}, tt.repoErr
				},
			}

			uc := NewNoteUsecase(mock)

			_, err := uc.GetByID(tt.id)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		note    domain.Note
		repoErr error
		wantErr error
	}{
		{
			name: "success",
			id:   "1",
			note: domain.Note{
				Title: "Test",
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   "1",
			note: domain.Note{
				Title: "Test",
			},
			repoErr: domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
		{
			name: "repo failure",
			id:   "1",
			note: domain.Note{
				Title: "Test",
			},
			repoErr: domain.ErrDb,
			wantErr: domain.ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := &mockRepo{
				updateFn: func(id string, note domain.Note) error {
					return tt.repoErr
				},
			}

			uc := NewNoteUsecase(mock)

			err := uc.Update(tt.id, tt.note)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		repoErr error
		wantErr error
	}{
		{
			name: "success",
			id:   "1",
		},
		{
			name:    "not found",
			id:      "1",
			repoErr: domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := &mockRepo{
				deleteFn: func(id string) error {
					return tt.repoErr
				},
			}

			uc := NewNoteUsecase(mock)

			err := uc.Delete(tt.id)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
