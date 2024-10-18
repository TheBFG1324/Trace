package executor

import (
	"sync"
	"trace/package/agent"
	"trace/package/task"
)

func ExecuteTask(a *agent.BaseAgent, t *task.Task) string {
	t.UpdateStatus(task.InProgress)
	t.UpdateOwner(a.GetID())
	t.DisplayTask()

	var wg sync.WaitGroup
	wg.Add(1)

	response := SimulateAPICall(a, t, &wg)
	wg.Wait()
	t.UpdateResult(response)
	t.UpdateStatus(task.Finished)
	return response
}

func SimulateAPICall(a *agent.BaseAgent, t *task.Task, wg *sync.WaitGroup) string {
	defer wg.Done()
	return "simulated response"
}
