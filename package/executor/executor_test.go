package executor_test

import (
	"testing"
	"trace/package/agent"
	"trace/package/executor"
	"trace/package/task"
)

// TestExecuteTask tests the ExecuteTask function.
func TestExecuteTask(t *testing.T) {
	// Define the JSON template for the agent
	jsonTemplate := map[string]interface{}{
		"taskName": "[[taskName]]",
		"age":      "[[age]]",
		"config": map[string]interface{}{
			"timeout": "[[timeout]]",
		},
	}

	// Create a new agent
	agent := agent.NewBaseAgent(
		"agent1",
		"Example Agent",
		"TypeA",
		"http://example.com/api",
		jsonTemplate,
		[]string{"capability1"},
	)

	// Define task parameters and global data
	taskParameters := map[string]interface{}{
		"taskName": "Process Data",
		"age":      6,
	}

	globalData := map[string]interface{}{
		"timeout": 30,
	}

	// Create a new task
	task := task.CreateTask(1, "Example Task", taskParameters)

	// Execute the task
	response := executor.ExecuteTask(agent, task, globalData)

	// Verify the response
	expected := "simulated response" // This should be the expected response from SimulateAPICall.
	if response != expected {
		t.Errorf("Expected response %s, but got %s", expected, response)
	}
}
