package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"todo-cli/internal/repository"
	"todo-cli/internal/usecase"
)

var taskUsecase *usecase.TaskUsecase

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple CLI Todo application to manage your tasks",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	repo := &repository.JSONRepository{Filename: "tasks.json"}
	taskUsecase = usecase.NewTaskUsecase(repo)
}
