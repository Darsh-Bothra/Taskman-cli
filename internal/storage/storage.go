// Package storage handles reading and writing task data to disk.
package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"todo-cli.com/internal/task"
)

// Store persists a Manager's state to a JSON file at path.
func Store(path string, m *task.Manager) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// Load reads JSON from path and populates the manager.
// If the file does not exist the manager is left in its zero state —
// this is not treated as an error (first run).
func Load(path string, m *task.Manager) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil // first run — no file yet
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}