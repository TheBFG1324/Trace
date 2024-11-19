package executor

import (
    "fmt"
    "sync"
    "time"
    "trace/package/agent"
    "trace/package/parser"
    "trace/package/task"
    "trace/package/utils/template"
)

// ExecuteTask performs the task using the provided agent and updates the task status accordingly.
func ExecuteTask(a *agent.BaseAgent, t *task.Task, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission) error {
    t.UpdateStatus(task.InProgress)
    t.UpdateOwner(a.GetID())
    t.DisplayTask()

    filteredGlobalData := FilterGlobalDataByPermissions(a.GetName(), globalPermissions, globalData)

    jsonPayload, err := template.LoadJSON(a.GetJsonBody(), t.Parameters, filteredGlobalData)
    if err != nil {
        return fmt.Errorf("error generating JSON payload: %w", err)
    }

    var wg sync.WaitGroup
    var response string

    wg.Add(1)
    go func() {
        defer wg.Done()
        response = SimulateAPICall(a, jsonPayload)
    }()
    wg.Wait()

    err = HandleResponse(a, t, globalData, globalPermissions, response)
    if err != nil {
        return fmt.Errorf("error handling response: %w", err)
    }

    t.UpdateStatus(task.Finished)
    return nil
}


// GetAgentPermissions retrieves the data permissions for a specific agent.
func GetAgentPermissions(a *agent.BaseAgent, globalPermissions map[string]*parser.Permission) map[string][]string {
    agentPermissions, ok := globalPermissions[a.GetName()]
    if !ok {
        return nil
    }
    return agentPermissions.DataPermissions
}

// FilterGlobalDataByPermissions creates a data mapping an agent can use based on its permissions.
func FilterGlobalDataByPermissions(agentName string, globalPermissions map[string]*parser.Permission, globalData map[string]*parser.Data) map[string]interface{} {
    agent := agent.SimulateLoadAgent("Name", agentName)
    if agent == nil {
        return nil
    }

    agentPermissions := GetAgentPermissions(agent, globalPermissions)
    if agentPermissions == nil {
        return nil
    }

    dependentGlobalData := make(map[string]interface{})

    for variable, permissions := range agentPermissions {
        if HasPermission(permissions, "READ") {
            if value, found := globalData[variable]; found {
                dependentGlobalData[variable] = value
            }
        }
    }

    return dependentGlobalData
}

// HasPermission checks if a specific permission exists in the list.
func HasPermission(permissions []string, target string) bool {
    for _, permission := range permissions {
        if permission == target {
            return true
        }
    }
    return false
}

// HandleResponse updates the global data if commanded
func HandleResponse(a *agent.BaseAgent, t *task.Task, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission, response string) error {
    variableRaw, hasOutputParameter := t.Parameters["OUTPUT"]
    if !hasOutputParameter {
        t.UpdateResult(response)
        return nil
    }

    variable, ok := variableRaw.(string)
    if !ok {
        return fmt.Errorf("expected OUTPUT parameter to be a string, got %T", variableRaw)
    }

    agentPermissions := GetAgentPermissions(a, globalPermissions)
    if agentPermissions == nil {
        return fmt.Errorf("agent '%s' does not have any permissions defined", a.GetName())
    }

    permissions, variableExists := agentPermissions[variable]
    if !variableExists || !HasPermission(permissions, "WRITE") {
        return fmt.Errorf("agent '%s' does not have WRITE permission for variable '%s'", a.GetName(), variable)
    }

    data, found := globalData[variable]
    if !found {
        return fmt.Errorf("variable '%s' not found in global data", variable)
    }

    data.Mu.Lock()
    data.InitialValue = response
    data.Mu.Unlock()

    newParam := make(map[string]interface{})
    newParam["OUTPUT"] = response
    t.UpdateParameters(newParam)

    t.UpdateResult(response)
    return nil
}


// SimulateAPICall simulates sending a payload to the agent's endpoint.
func SimulateAPICall(a *agent.BaseAgent, jsonPayload string) string {
    fmt.Println("Sending payload to", a.GetEndpoint())
    fmt.Println("Payload:", jsonPayload)
    time.Sleep(2 * time.Second)
    return "simulated response"
}
