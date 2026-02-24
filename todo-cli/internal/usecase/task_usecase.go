package usecase

import (
	"errors"
	"todo-cli/internal/domain"
	"todo-cli/internal/repository"
)

type TaskUsecase struct {
	repo repository.TaskRepository
}

func NewTaskUsecase(r repository.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: r}
}

func (u *TaskUsecase) Add(name string) error {
	tasks, err := u.repo.Load()
	if err != nil {
		return err
	}

	id := 1
	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	task := domain.Task{
		ID:   id,
		Name: name,
		Done: false,
	}

	tasks = append(tasks, task)
	return u.repo.Save(tasks)
}

func (u *TaskUsecase) List() ([]domain.Task, error) {
	return u.repo.Load()
}

func (u *TaskUsecase) Delete(id int) error {
	tasks, err := u.repo.Load()
	if err != nil {
		return err
	}

	var updated []domain.Task
	found := false

	for _, task := range tasks {
		if task.ID != id {
			updated = append(updated, task)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("TASK NOT FOUND")
	}

	// Renumber sequentially
	for i := range updated {
		updated[i].ID = i + 1
	}

	return u.repo.Save(updated)
}

func (u *TaskUsecase) MarkDone(id int) error {
	tasks, err := u.repo.Load()
	if err != nil {
		return err
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}

	if !found {
		return errors.New("TASK NOT FOUND")
	}

	return u.repo.Save(tasks)
}
