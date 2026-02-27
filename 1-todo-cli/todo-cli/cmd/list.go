package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := taskUsecase.List()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, t := range tasks {
			status := " "
			if t.Done {
				status = "✓"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
