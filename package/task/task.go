// pkg/task/task.go
package task

import (
	"fmt"
	"sync"
)

// Status represents the status of a task.
type Status int

const (
	Pending Status = iota
	Claimed
	InProgress
	Finished
)

// Task represents a unit of work.
type Task struct {
	ID          int
	Description string
	Owner       string
	Status      Status
	Parameters  map[string]interface{}
	Result      []string
	mu sync.Mutex
}

// CreateTask initializes a new task.
func CreateTask(id int, description string, parameters map[string]interface{}) *Task {
	return &Task{
		ID:          id,
		Description: description,
		Owner:       "None",
		Status:      Pending,
		Parameters:  parameters,
		Result:      []string{},
	}
}

// UpdateStatus updates the task's status.
func (t *Task) UpdateStatus(newStatus Status) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Status = newStatus
}

// UpdateOwner updates the task's owner
func (t *Task) UpdateOwner(ownerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	t.Owner = ownerID
}

// UpdateParameters merges new parameters into the task's parameters.
func (t *Task) UpdateParameters(newParameters map[string]interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for key, value := range newParameters {
		t.Parameters[key] = value
	}
}

// UpdateResult appends an item to the task's results.
func (t *Task) UpdateResult(item string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Result = append(t.Result, item)
}

// DisplayTask prints the task's details.
func (t *Task) GetInfoString() string {
    return fmt.Sprintf(
        "Task ID: %d\nDescription: %s\nStatus: %d\nOwner: %v\nParameters: %v\nResults: %v\n",
        t.ID, t.Description, t.Status, t.Owner, t.Parameters, t.Result,
    )
}
