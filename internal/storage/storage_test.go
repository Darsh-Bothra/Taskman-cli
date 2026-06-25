package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"todo-cli.com/internal/storage"
	"todo-cli.com/internal/task"
)

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	// Populate a manager and persist it.
	m1 := task.New()
	m1.Add("first task")
	m1.Add("second task")
	_ = m1.Complete(1)

	if err := storage.Store(path, m1); err != nil {
		t.Fatalf("Store: %v", err)
	}

	// Load into a fresh manager and compare.
	m2 := task.New()
	if err := storage.Load(path, m2); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(m2.Tasks) != 2 {
		t.Errorf("Tasks len = %d, want 2", len(m2.Tasks))
	}
	if !m2.Tasks[0].Completed {
		t.Error("first task should be completed after reload")
	}
	if m2.NextID != m1.NextID {
		t.Errorf("NextID mismatch: got %d, want %d", m2.NextID, m1.NextID)
	}
}

func TestLoadMissingFile(t *testing.T) {
	m := task.New()
	err := storage.Load("/nonexistent/path/tasks.json", m)
	if err != nil {
		t.Errorf("Load on missing file should return nil, got: %v", err)
	}
}

func TestStoreCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new_tasks.json")

	m := task.New()
	m.Add("test")

	if err := storage.Store(path, m); err != nil {
		t.Fatalf("Store: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Store should have created the file")
	}
}