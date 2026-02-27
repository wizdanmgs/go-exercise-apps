package usecase

import (
	"errors"
	"reflect"
	"testing"
	"todo-cli/internal/domain"
)

type mockRepository struct {
	tasks []domain.Task
	err   error
}

func (m *mockRepository) Load() ([]domain.Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tasks, nil
}

func (m *mockRepository) Save(tasks []domain.Task) error {
	if m.err != nil {
		return m.err
	}
	m.tasks = tasks
	return nil
}

func TestAdd(t *testing.T) {
	mockRepo := &mockRepository{
		tasks: []domain.Task{},
	}

	u := NewTaskUsecase(mockRepo)

	err := u.Add("Learn Testing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockRepo.tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(mockRepo.tasks))
	}

	if mockRepo.tasks[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", mockRepo.tasks[0].ID)
	}

	if mockRepo.tasks[0].Name != "Learn Testing" {
		t.Fatalf("unexpected task name: %s", mockRepo.tasks[0].Name)
	}
}

func TestAdd_RepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		err: errors.New("db failure"),
	}

	u := NewTaskUsecase(mockRepo)

	err := u.Add("X")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_Renumber(t *testing.T) {
	mockRepo := &mockRepository{
		tasks: []domain.Task{
			{ID: 1, Name: "A"},
			{ID: 2, Name: "B"},
			{ID: 3, Name: "C"},
		},
	}

	u := NewTaskUsecase(mockRepo)

	err := u.Delete(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []domain.Task{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "C"},
	}

	if !reflect.DeepEqual(mockRepo.tasks, expected) {
		t.Fatalf("expected tasks to be: %v, got: %v", expected, mockRepo.tasks)
	}
}

func TestDelete_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		tasks: []domain.Task{
			{ID: 1, Name: "A"},
		},
	}

	u := NewTaskUsecase(mockRepo)

	err := u.Delete(99)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMarkDone(t *testing.T) {
	mockRepo := &mockRepository{
		tasks: []domain.Task{
			{ID: 1, Name: "A", Done: false},
		},
	}

	u := NewTaskUsecase(mockRepo)

	err := u.MarkDone(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockRepo.tasks[0].Done {
		t.Fatal("expected task to be marked as done")
	}
}
