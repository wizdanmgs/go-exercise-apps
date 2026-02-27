package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete task by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}

		err = taskUsecase.Delete(id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Task deleted!")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
