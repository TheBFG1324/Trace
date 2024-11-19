package executor

import (
    "sync"
    "testing"
    "trace/package/agent"
    "trace/package/parser"
    "trace/package/task"
)

// TestExecuteTask_Success tests the ExecuteTask function with valid inputs.
func TestExecuteTask_Success(t *testing.T) {
    agent := agent.SimulateLoadAgent("Name", "Flight Getter")
    if agent == nil {
        t.Fatal("Agent not found")
    }

    task := task.CreateTask(1, "Book Flight", map[string]interface{}{
        "origin":      "NYC",
        "destination": "LAX",
        "date":        "2023-10-10",
        "OUTPUT":      "flightInfo",
    })

    globalData := map[string]*parser.Data{
        "flightInfo": {
            DataName:     "flightInfo",
            DataType:     "String",
            InitialValue: "",
            Mu:           sync.Mutex{},
        },
    }

    globalPermissions := map[string]*parser.Permission{
        "Flight Getter": {
            AgentName: "Flight Getter",
            DataPermissions: map[string][]string{
                "flightInfo": {"READ", "WRITE"},
            },
        },
    }

    err := ExecuteTask(agent, task, globalData, globalPermissions)
    if err != nil {
        t.Fatalf("ExecuteTask failed: %v", err)
    }

    expectedValue := "simulated response"
    actualValue := globalData["flightInfo"].InitialValue
    if actualValue != expectedValue {
        t.Errorf("Expected globalData['flightInfo'].InitialValue to be '%s', got '%s'", expectedValue, actualValue)
    }
}

// TestExecuteTask_NoWritePermission tests ExecuteTask when the agent lacks WRITE permission.
func TestExecuteTask_NoWritePermission(t *testing.T) {
    agent := agent.SimulateLoadAgent("Name", "Flight Getter")
    if agent == nil {
        t.Fatal("Agent not found")
    }

    task := task.CreateTask(1, "Book Flight", map[string]interface{}{
        "origin":      "NYC",
        "destination": "LAX",
        "date":        "2023-10-10",
        "OUTPUT":      "flightInfo",
    })

    globalData := map[string]*parser.Data{
        "flightInfo": {
            DataName:     "flightInfo",
            DataType:     "String",
            InitialValue: "",
            Mu:           sync.Mutex{},
        },
    }

    globalPermissions := map[string]*parser.Permission{
        "Flight Getter": {
            AgentName: "Flight Getter",
            DataPermissions: map[string][]string{
                "flightInfo": {"READ"}, 
            },
        },
    }

    err := ExecuteTask(agent, task, globalData, globalPermissions)
    if err == nil {
        t.Fatal("Expected error due to lack of WRITE permission, but got none")
    }
}

func TestExecuteTask_MissingGlobalData(t *testing.T) {
    agent := agent.SimulateLoadAgent("Name", "Flight Getter")
    if agent == nil {
        t.Fatal("Agent not found")
    }

    task := task.CreateTask(1, "Book Flight", map[string]interface{}{
        "origin":      "NYC",
        "destination": "LAX",
        "date":        "2023-10-10",
        "OUTPUT":      "flightInfo",
    })

    globalData := map[string]*parser.Data{}

    globalPermissions := map[string]*parser.Permission{
        "Flight Getter": {
            AgentName: "Flight Getter",
            DataPermissions: map[string][]string{
                "flightInfo": {"READ", "WRITE"},
            },
        },
    }

    err := ExecuteTask(agent, task, globalData, globalPermissions)
    if err == nil {
        t.Fatal("Expected error due to missing global data, but got none")
    }
}
