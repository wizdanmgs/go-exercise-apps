package repository

import (
	"encoding/json"
	"os"
	"todo-cli/internal/domain"
)

type JSONRepository struct {
	Filename string
}

func (r *JSONRepository) Load() ([]domain.Task, error) {
	var tasks []domain.Task

	data, err := os.ReadFile(r.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil // return empty slice if file not exist
		}
		return nil, err
	}

	err = json.Unmarshal(data, &tasks)

	return tasks, err
}

func (r *JSONRepository) Save(tasks []domain.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.Filename, data, 0644)
}
