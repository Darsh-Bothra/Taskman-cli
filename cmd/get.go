package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"todo-cli.com/internal/task"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Show a single task by ID",
	Long:  "Print the full details of one task.",
	Example: `  taskman get 3`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func runGet(_ *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID %q — must be a number", args[0])
	}

	t := manager.GetByID(id)
	if t == nil {
		return fmt.Errorf("%w: #%d", task.ErrNotFound, id)
	}

	status := "Pending"
	if t.Completed {
		status = "Done"
	}

	fmt.Printf("ID:          %d\n", t.ID)
	fmt.Printf("Description: %s\n", t.Description)
	fmt.Printf("Status:      %s\n", status)
	fmt.Printf("Created:     %s\n", t.CreatedAt.Format("2006-01-02 15:04:05"))
	if t.CompletedAt != nil {
		fmt.Printf("Completed:   %s\n", t.CompletedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

// Ensure ErrNotFound is in scope for error wrapping reference.
var _ = errors.Is