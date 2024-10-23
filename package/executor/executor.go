// trace/package/executor/executor.go
package executor

import (
    "fmt"
    "sync"
    "trace/package/agent"
    "trace/package/task"
    "trace/package/utils/template"
	"time"
)

// ExecuteTask performs the task using the provided agent and updates the task status accordingly.
func ExecuteTask(a *agent.BaseAgent, t *task.Task, globalData map[string]interface{}) string {
    t.UpdateStatus(task.InProgress)
    t.UpdateOwner(a.GetID())
    t.DisplayTask()

    // Generate the JSON payload using the template package
    jsonPayload, err := template.LoadJSON(a.GetJsonBody(), t.Parameters, globalData)
    if err != nil {
        fmt.Println("Error generating JSON payload:", err)
        return ""
    }

    var wg sync.WaitGroup
    var response string

    wg.Add(1)
    go func() {
        defer wg.Done()
        response = SimulateAPICall(a, jsonPayload)
    }()

    wg.Wait()
    t.UpdateResult(response)
    t.UpdateStatus(task.Finished)
    return response
}

// SimulateAPICall simulates sending a payload to the agent's endpoint.
func SimulateAPICall(a *agent.BaseAgent, jsonPayload string) string {
    fmt.Println("Sending payload to", a.GetEndpoint())
    fmt.Println("Payload:", jsonPayload)
    time.Sleep(2 * time.Second)
    return "simulated response"
}
