package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark task as done",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}

		err = taskUsecase.MarkDone(id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Task marked as done!")
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
