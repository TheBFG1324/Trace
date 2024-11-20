package task_test

import (
	"fmt"
	"testing"
	"trace/package/task"
)

// TestCreateTask tests the creation of a new task.
func TestCreateTask(t *testing.T) {
	params := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	taskInstance := task.CreateTask(1, "Test Task", params)

	if taskInstance.ID != 1 {
		t.Errorf("Expected task ID to be 1, but got %d", taskInstance.ID)
	}
	if taskInstance.Description != "Test Task" {
		t.Errorf("Expected description to be 'Test Task', but got %s", taskInstance.Description)
	}
	if taskInstance.Status != task.Pending {
		t.Errorf("Expected task status to be Pending, but got %d", taskInstance.Status)
	}
}

// TestUpdateStatus tests updating the task's status.
func TestUpdateStatus(t *testing.T) {
	taskInstance := task.CreateTask(2, "Update Status Task", nil)
	taskInstance.UpdateStatus(task.InProgress)

	if taskInstance.Status != task.InProgress {
		t.Errorf("Expected task status to be InProgress, but got %d", taskInstance.Status)
	}
}

// TestUpdateOwner tests updating the task's owner.
func TestUpdateOwner(t *testing.T) {
	taskInstance := task.CreateTask(3, "Update Owner Task", nil)
	taskInstance.UpdateOwner("user123")

	if taskInstance.Owner != "user123" {
		t.Errorf("Expected owner to be 'user123', but got %s", taskInstance.Owner)
	}
}

// TestUpdateParameters tests updating the task's parameters.
func TestUpdateParameters(t *testing.T) {
	taskInstance := task.CreateTask(4, "Update Parameters Task", map[string]interface{}{
		"key1": "value1",
	})

	newParams := map[string]interface{}{
		"key2": "value2",
		"key1": "newValue1",
	}
	taskInstance.UpdateParameters(newParams)

	if taskInstance.Parameters["key1"] != "newValue1" {
		t.Errorf("Expected key1 to be 'newValue1', but got %v", taskInstance.Parameters["key1"])
	}
	if taskInstance.Parameters["key2"] != "value2" {
		t.Errorf("Expected key2 to be 'value2', but got %v", taskInstance.Parameters["key2"])
	}
}

// TestUpdateResult tests updating the task's results.
func TestUpdateResult(t *testing.T) {
	taskInstance := task.CreateTask(5, "Update Result Task", nil)
	taskInstance.UpdateResult("result1")
	taskInstance.UpdateResult("result2")

	if len(taskInstance.Result) != 2 {
		t.Errorf("Expected 2 results, but got %d", len(taskInstance.Result))
	}
	if taskInstance.Result[0] != "result1" {
		t.Errorf("Expected first result to be 'result1', but got %s", taskInstance.Result[0])
	}
	if taskInstance.Result[1] != "result2" {
		t.Errorf("Expected second result to be 'result2', but got %s", taskInstance.Result[1])
	}
}

// TestGetInfoString tests displaying the task details.
func TestGetInfoString(t *testing.T) {
	taskInstance := task.CreateTask(6, "Display Task", map[string]interface{}{
		"key1": "value1",
	})
	taskInstance.UpdateOwner("user456")
	taskInstance.UpdateStatus(task.Finished)
	taskInstance.UpdateResult("result1")

	info := taskInstance.GetInfoString()
	fmt.Println(info)
}
