package executor_test

import (
	"sync"
	"testing"
	"trace/package/agent"
	"trace/package/executor"
	"trace/package/logger"
	"trace/package/parser"
)

// TestExecuteTask_Success verifies the successful execution of a task.
func TestExecuteTask_Success(t *testing.T) {
	// Load a mock agent
	mockAgent := agent.SimulateLoadAgent("Name", "Flight Getter")
	if mockAgent == nil {
		t.Fatal("Agent not found")
	}

	// Create a parser task with valid parameters
	mockTask := &parser.Task{
		TaskName:  "Book Flight",
		AgentName: "Flight Getter",
		Parameters: map[string]string{
			"origin":      "NYC",
			"destination": "LAX",
			"date":        "2023-10-10",
			"OUTPUT":      "flightInfo",
		},
	}

	// Define global data with synchronization
	globalData := map[string]*parser.Data{
		"flightInfo": {
			DataName:     "flightInfo",
			DataType:     "String",
			InitialValue: "",
			Mu:           sync.Mutex{},
		},
	}

	// Define global permissions for the mock agent
	globalPermissions := map[string]*parser.Permission{
		"Flight Getter": {
			AgentName: "Flight Getter",
			DataPermissions: map[string][]string{
				"flightInfo": {"READ", "WRITE"},
			},
		},
	}

	// Create a logger instance
	log := logger.NewLogger()

	// Execute the task
	err := executor.ExecuteTask(mockAgent.GetName(), mockTask, globalData, globalPermissions, log)
	if err != nil {
		t.Fatalf("ExecuteTask failed: %v", err)
	}

	// Check if the global data was updated as expected
	expectedValue := "simulated response"
	actualValue := globalData["flightInfo"].InitialValue
	if actualValue != expectedValue {
		t.Errorf("Expected globalData['flightInfo'].InitialValue to be '%s', got '%s'", expectedValue, actualValue)
	}

	// Print logs for debugging purposes
	log.PrintAllLogs()
}

// TestExecuteTask_NoWritePermission verifies behavior when the agent lacks WRITE permission.
func TestExecuteTask_NoWritePermission(t *testing.T) {
	// Load a mock agent
	mockAgent := agent.SimulateLoadAgent("Name", "Flight Getter")
	if mockAgent == nil {
		t.Fatal("Agent not found")
	}

	// Create a parser task
	mockTask := &parser.Task{
		TaskName:  "Book Flight",
		AgentName: "Flight Getter",
		Parameters: map[string]string{
			"origin":      "NYC",
			"destination": "LAX",
			"date":        "2023-10-10",
			"OUTPUT":      "flightInfo",
		},
	}

	// Define global data
	globalData := map[string]*parser.Data{
		"flightInfo": {
			DataName:     "flightInfo",
			DataType:     "String",
			InitialValue: "",
			Mu:           sync.Mutex{},
		},
	}

	// Define global permissions without WRITE permission
	globalPermissions := map[string]*parser.Permission{
		"Flight Getter": {
			AgentName: "Flight Getter",
			DataPermissions: map[string][]string{
				"flightInfo": {"READ"}, // Lacks WRITE permission
			},
		},
	}

	// Create a logger instance
	log := logger.NewLogger()

	// Execute the task
	err := executor.ExecuteTask(mockAgent.GetName(), mockTask, globalData, globalPermissions, log)
	if err == nil {
		t.Fatal("Expected error due to lack of WRITE permission, but got none")
	}

	// Print logs for debugging purposes
	log.PrintAllLogs()
}

// TestExecuteTask_MissingGlobalData verifies behavior when required global data is missing.
func TestExecuteTask_MissingGlobalData(t *testing.T) {
	// Load a mock agent
	mockAgent := agent.SimulateLoadAgent("Name", "Flight Getter")
	if mockAgent == nil {
		t.Fatal("Agent not found")
	}

	// Create a parser task
	mockTask := &parser.Task{
		TaskName:  "Book Flight",
		AgentName: "Flight Getter",
		Parameters: map[string]string{
			"origin":      "NYC",
			"destination": "LAX",
			"date":        "2023-10-10",
			"OUTPUT":      "flightInfo",
		},
	}

	// Define empty global data to simulate missing data
	globalData := map[string]*parser.Data{}

	// Define global permissions
	globalPermissions := map[string]*parser.Permission{
		"Flight Getter": {
			AgentName: "Flight Getter",
			DataPermissions: map[string][]string{
				"flightInfo": {"READ", "WRITE"},
			},
		},
	}

	// Create a logger instance
	log := logger.NewLogger()

	// Execute the task
	err := executor.ExecuteTask(mockAgent.GetName(), mockTask, globalData, globalPermissions, log)
	if err == nil {
		t.Fatal("Expected error due to missing global data, but got none")
	}

	// Print logs for debugging purposes
	log.PrintAllLogs()
}
