package task

import (
	"errors"
	"sync"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Manager struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		tasks: make(map[string]*Task),
	}
}

// Create a new task and start execution
func (m *Manager) CreateTask() *Task {
	task := NewTask()

	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks[task.ID] = task

	return task
}

// Get returns a task by ID
func (m *Manager) Get(id string) (*Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// Delete cancels and removes the task
func (m *Manager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}

	task.Cancel()
	delete(m.tasks, id)
	return nil
}
