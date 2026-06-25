package task_test

import (
	"testing"

	"todo-cli.com/internal/task"
)

func TestAdd(t *testing.T) {
	m := task.New()

	got := m.Add("buy milk")
	if got.ID != 1 {
		t.Errorf("first task ID = %d, want 1", got.ID)
	}
	if got.Description != "buy milk" {
		t.Errorf("description = %q, want %q", got.Description, "buy milk")
	}
	if got.Completed {
		t.Error("new task should not be completed")
	}
	if m.NextID != 2 {
		t.Errorf("NextID = %d, want 2", m.NextID)
	}
}

func TestGetByID(t *testing.T) {
	m := task.New()
	m.Add("task one")
	m.Add("task two")

	got := m.GetByID(2)
	if got == nil {
		t.Fatal("GetByID(2) returned nil, want task")
	}
	if got.Description != "task two" {
		t.Errorf("description = %q, want %q", got.Description, "task two")
	}

	if m.GetByID(99) != nil {
		t.Error("GetByID(99) should return nil for unknown ID")
	}
}

func TestComplete(t *testing.T) {
	m := task.New()
	m.Add("write tests")

	if err := m.Complete(1); err != nil {
		t.Fatalf("Complete(1) error: %v", err)
	}

	got := m.GetByID(1)
	if !got.Completed {
		t.Error("task should be marked completed")
	}
	if got.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	// Completing again should return ErrAlreadyComplete.
	if err := m.Complete(1); err != task.ErrAlreadyComplete {
		t.Errorf("second Complete got %v, want ErrAlreadyComplete", err)
	}

	// Unknown ID.
	if err := m.Complete(99); err != task.ErrNotFound {
		t.Errorf("Complete(99) got %v, want ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	m := task.New()
	m.Add("temporary task")

	if err := m.Delete(1); err != nil {
		t.Fatalf("Delete(1) error: %v", err)
	}
	if len(m.Tasks) != 0 {
		t.Errorf("Tasks len = %d, want 0", len(m.Tasks))
	}

	// Deleting again should return ErrNotFound.
	if err := m.Delete(1); err != task.ErrNotFound {
		t.Errorf("second Delete got %v, want ErrNotFound", err)
	}
}

func TestPending(t *testing.T) {
	m := task.New()
	m.Add("pending task")
	m.Add("done task")
	_ = m.Complete(2)

	pending := m.Pending()
	if len(pending) != 1 {
		t.Errorf("Pending() len = %d, want 1", len(pending))
	}
	if pending[0].ID != 1 {
		t.Errorf("Pending()[0].ID = %d, want 1", pending[0].ID)
	}
}

func TestAll(t *testing.T) {
	m := task.New()
	m.Add("a")
	m.Add("b")

	all := m.All()
	if len(all) != 2 {
		t.Errorf("All() len = %d, want 2", len(all))
	}

	// Verify All() returns a copy — mutations should not affect manager.
	all[0].Description = "mutated"
	if m.Tasks[0].Description == "mutated" {
		t.Error("All() should return a copy, not a reference")
	}
}