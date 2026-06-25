package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"todo-cli.com/internal/task"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Mark a task as done",
	Long:  "Mark the task with the given ID as completed.",
	Example: `  taskman complete 2`,
	Args: cobra.ExactArgs(1),
	RunE: runComplete,
}

func runComplete(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID %q — must be a number", args[0])
	}

	if err := manager.Complete(id); err != nil {
		switch {
		case errors.Is(err, task.ErrNotFound):
			return fmt.Errorf("task #%d not found", id)
		case errors.Is(err, task.ErrAlreadyComplete):
			return fmt.Errorf("task #%d is already completed", id)
		}
		return err
	}

	// Re-fetch to display the description.
	t := manager.GetByID(id)
	fmt.Printf("✓ Completed task #%d: %s\n", id, t.Description)
	return nil
}