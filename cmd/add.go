package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add a new task",
	Long:  "Add a new task with a description. Wrap multi-word descriptions in quotes.",
	Example: `  taskman add "Buy groceries"
  taskman add Write the quarterly report`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

func runAdd(_ *cobra.Command, args []string) error {
	description := strings.Join(args, " ")
	t := manager.Add(description)
	fmt.Printf("✓ Added task #%d: %s\n", t.ID, t.Description)
	return nil
}