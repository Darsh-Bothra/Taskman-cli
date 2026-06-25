package cmd

import (
	"fmt"
	"strings"

	"todo-cli.com/internal/task"
	"github.com/spf13/cobra"
)

var showAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  "List pending tasks. Pass --all to include completed tasks.",
	Example: `  taskman list
  taskman list --all`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "include completed tasks")
}

func runList(_ *cobra.Command, _ []string) error {
	var tasks []task.Task
	if showAll {
		tasks = manager.All()
	} else {
		tasks = manager.Pending()
	}

	if len(tasks) == 0 {
		if showAll {
			fmt.Println("No tasks found. Add one with: taskman add <description>")
		} else {
			fmt.Println("No pending tasks. Run 'taskman list --all' to see completed tasks.")
		}
		return nil
	}

	fmt.Printf("%-5s %-12s %-50s %s\n", "ID", "STATUS", "DESCRIPTION", "CREATED")
	fmt.Println(strings.Repeat("─", 85))

	for _, t := range tasks {
		status := "⏳ Pending"
		if t.Completed {
			status = "✓ Done"
		}
		fmt.Printf("%-5d %-12s %-50s %s\n",
			t.ID,
			status,
			truncate(t.Description, 49),
			t.CreatedAt.Format("2006-01-02"),
		)
	}
	return nil
}

// truncate shortens s to max runes, adding "…" if needed.
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}