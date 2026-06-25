package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"todo-cli.com/internal/task"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a task by ID",
	Long:  "Permanently remove the task with the given ID.",
	Example: `  taskman delete 5`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func runDelete(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID %q — must be a number", args[0])
	}

	// Capture the description before deletion for the success message.
	t := manager.GetByID(id)
	if t == nil {
		return fmt.Errorf("%w: #%d", task.ErrNotFound, id)
	}
	description := t.Description

	if err := manager.Delete(id); err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return fmt.Errorf("task #%d not found", id)
		}
		return err
	}

	fmt.Printf("🗑  Deleted task #%d: %s\n", id, description)
	return nil
}