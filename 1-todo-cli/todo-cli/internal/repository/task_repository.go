package repository

import "todo-cli/internal/domain"

type TaskRepository interface {
	Load() ([]domain.Task, error)
	Save([]domain.Task) error
}
