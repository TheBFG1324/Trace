package executor

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"trace/package/agent"
	"trace/package/logger"
	"trace/package/parser"
	"trace/package/task"
	"trace/package/utils/template"
)

// ExecuteTask performs the task using the provided agent and updates the task status accordingly.
func ExecuteTask(a *agent.BaseAgent, t *task.Task, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission, l *logger.Logger) error {
	var logs []logger.Log
	var wg sync.WaitGroup
	var mu sync.Mutex

	t.UpdateStatus(task.InProgress)
	t.UpdateOwner(a.GetID())
	logs = append(logs, logger.NewLog("Starting Task: "+t.GetInfoString()))

	filteredGlobalData := FilterGlobalDataByPermissions(a.GetName(), globalPermissions, globalData)
	filteredDataStr, err := json.Marshal(filteredGlobalData)
	if err != nil {
		logs = append(logs, logger.NewLog("Error marshalling filtered global data: "+err.Error()))
	} else {
		logs = append(logs, logger.NewLog("Filtered global data for agent "+a.GetName()+": "+string(filteredDataStr)))
	}

	jsonPayload, err := template.LoadJSON(a.GetJsonBody(), t.Parameters, filteredGlobalData)
	if err != nil {
		logs = append(logs, logger.NewLog("Error generating JSON payload: "+err.Error()))
		l.AddLogs(logs)
		return fmt.Errorf("error generating JSON payload: %w", err)
	}
	logs = append(logs, logger.NewLog("JSON Payload: "+jsonPayload))

	var response string
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp := SimulateAPICall(a, jsonPayload)
		mu.Lock()
		response = resp
		mu.Unlock()
	}()

	wg.Wait()
	logs = append(logs, logger.NewLog("Response from endpoint: "+response))

	err = HandleResponse(a, t, globalData, globalPermissions, response)
	if err != nil {
		logs = append(logs, logger.NewLog("Error handling response: "+err.Error()))
		l.AddLogs(logs)
		return fmt.Errorf("error handling response: %w", err)
	}

	globalDataStr := GlobalDataToString(globalData)
	logs = append(logs, logger.NewLog("Updated Global Data: "+globalDataStr))

	t.UpdateStatus(task.Finished)
	logs = append(logs, logger.NewLog("Task Status: "+t.GetInfoString()))

	l.AddLogs(logs)

	return nil
}

// GlobalDataToString converts the global data map to a JSON string for logging.
func GlobalDataToString(globalData map[string]*parser.Data) string {
	dataCopy := make(map[string]interface{})

	for key, data := range globalData {
		data.Mu.Lock()
		dataCopy[key] = map[string]interface{}{
			"DataName":     data.DataName,
			"DataType":     data.DataType,
			"InitialValue": data.InitialValue,
		}
		data.Mu.Unlock()
	}

	bytes, err := json.Marshal(dataCopy)
	if err != nil {
		return "Error converting global data to string: " + err.Error()
	}
	return string(bytes)
}

// GetAgentPermissions retrieves the data permissions for a specific agent.
func GetAgentPermissions(a *agent.BaseAgent, globalPermissions map[string]*parser.Permission) map[string][]string {
	agentPermissions, ok := globalPermissions[a.GetName()]
	if !ok {
		return nil
	}
	return agentPermissions.DataPermissions
}

// FilterGlobalDataByPermissions filters global data based on agent's permissions.
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
				value.Mu.Lock()
				dependentGlobalData[variable] = value.InitialValue
				value.Mu.Unlock()
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

// HandleResponse updates the global data if required.
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

	t.UpdateResult(response)
	return nil
}

// SimulateAPICall simulates sending a payload to the agent's endpoint.
func SimulateAPICall(a *agent.BaseAgent, jsonPayload string) string {
	time.Sleep(2 * time.Second)
	return "simulated response"
}
