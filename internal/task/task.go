// Package task defines the core domain types for taskman.
package task

import (
	"errors"
	"time"
)

// Common sentinel errors returned by Manager operations.
var (
	ErrNotFound        = errors.New("task not found")
	ErrAlreadyComplete = errors.New("task is already completed")
)

// Task represents a single to-do item.
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Manager holds the in-memory list of tasks and is responsible for
// all business logic. Persistence is handled by the storage layer.
type Manager struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"next_id"`
}

// New returns a Manager ready to use (NextID starts at 1).
func New() *Manager {
	return &Manager{
		Tasks:  make([]Task, 0),
		NextID: 1,
	}
}

// Add creates a new Task and appends it to the manager.
// Returns the created Task.
func (m *Manager) Add(description string) Task {
	t := Task{
		ID:          m.NextID,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	m.Tasks = append(m.Tasks, t)
	m.NextID++
	return t
}

// GetByID returns a pointer to the task with the given ID, or nil if
// not found.  Callers that need a copy (not a pointer) should
// dereference after the nil check.
func (m *Manager) GetByID(id int) *Task {
	for i := range m.Tasks {
		if m.Tasks[i].ID == id {
			return &m.Tasks[i]
		}
	}
	return nil
}

// Delete removes the task with the given ID.
// Returns ErrNotFound if no task matches.
func (m *Manager) Delete(id int) error {
	for i, t := range m.Tasks {
		if t.ID == id {
			m.Tasks = append(m.Tasks[:i], m.Tasks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// Complete marks the task with the given ID as done.
// Returns ErrNotFound or ErrAlreadyComplete where appropriate.
func (m *Manager) Complete(id int) error {
	t := m.GetByID(id)
	if t == nil {
		return ErrNotFound
	}
	if t.Completed {
		return ErrAlreadyComplete
	}
	now := time.Now()
	t.Completed = true
	t.CompletedAt = &now
	return nil
}

// All returns a copy of the full task slice.
func (m *Manager) All() []Task {
	out := make([]Task, len(m.Tasks))
	copy(out, m.Tasks)
	return out
}

// Pending returns only tasks that are not yet completed.
func (m *Manager) Pending() []Task {
	var out []Task
	for _, t := range m.Tasks {
		if !t.Completed {
			out = append(out, t)
		}
	}
	return out
}