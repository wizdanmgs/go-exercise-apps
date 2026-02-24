package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := taskUsecase.Add(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Task added!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
