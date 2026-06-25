// Package cmd implements the taskman CLI commands.
package cmd

import (
	"fmt"
	"os"

	"todo-cli.com/internal/config"
	"todo-cli.com/internal/storage"
	"todo-cli.com/internal/task"
	"github.com/spf13/cobra"
)

// shared CLI state — intentionally package-scoped (not global vars) so
// sub-command files can read them without import cycles.
var (
	taskFile string         // --file flag value
	manager  *task.Manager  // loaded before every command, saved after
)

// rootCmd is the entry-point command (taskman).
var rootCmd = &cobra.Command{
	Use:   "taskman",
	Short: "A personal task manager",
	Long: `taskman — a fast, local CLI task manager.

Store tasks in a plain JSON file, mark them done, and
track your productivity without leaving the terminal.`,
	// Load tasks before every sub-command; save when it finishes.
	PersistentPreRunE:  loadTasks,
	PersistentPostRunE: saveTasks,
}

// Execute is the single public entry point called by main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	config.Init()

	rootCmd.PersistentFlags().StringVar(
		&taskFile, "file", "",
		"path to task file (default: $HOME/.taskman.json)",
	)

	rootCmd.AddCommand(addCmd, listCmd, getCmd, deleteCmd, completeCmd)
}

// loadTasks is run before every sub-command.
func loadTasks(_ *cobra.Command, _ []string) error {
	manager = task.New()
	path := config.TaskFile(taskFile)
	if err := storage.Load(path, manager); err != nil {
		return fmt.Errorf("loading tasks: %w", err)
	}
	return nil
}

// saveTasks is run after every sub-command (even if it fails, cobra still
// calls PostRunE — we only skip on PreRunE failure).
func saveTasks(_ *cobra.Command, _ []string) error {
	path := config.TaskFile(taskFile)
	if err := storage.Store(path, manager); err != nil {
		return fmt.Errorf("saving tasks: %w", err)
	}
	return nil
}