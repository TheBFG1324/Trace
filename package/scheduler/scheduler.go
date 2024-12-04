package scheduler

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"trace/package/executor"
	"trace/package/logger"
	"trace/package/parser"
)

// RunParentRequest schedules and runs the AICL parent request script
func RunParentRequest(p *parser.ParentRequest, l *logger.Logger) bool {
	errors := []string{}
	statements := p.Statements
	globalData := p.GlobalData
	globalPermissions := p.Permissions

	for _, stmt := range statements {
		RunStatement(stmt, globalData, globalPermissions, l, &errors)
	}

	if len(errors) != 0 {
		fmt.Println("Errors occurred during runtime:", errors)
		return false
	}
	return true
}

// RunStatement handles the execution of a single statement
func RunStatement(stmt interface{}, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission, l *logger.Logger, errors *[]string) {
	switch s := stmt.(type) {
	case *parser.Task:
		err := RunTask(s, globalData, globalPermissions, l)
		if err != nil {
			*errors = append(*errors, err.Error())
		}
	case *parser.RunSeqBlock:
		RunSeqBlock(s, globalData, globalPermissions, l, errors)
	case *parser.RunConBlock:
		RunConBlock(s, globalData, globalPermissions, l, errors)
	default:
		errMsg := "Unknown statement type"
		fmt.Println(errMsg)
		*errors = append(*errors, errMsg)
	}
}

// RunSeqBlock runs the tasks sequentially
func RunSeqBlock(seqBlock *parser.RunSeqBlock, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission, l *logger.Logger, errors *[]string) {
	for _, stmt := range seqBlock.Statements {
		RunStatement(stmt, globalData, globalPermissions, l, errors)
	}
}

// RunConBlock runs the tasks concurrently
func RunConBlock(conBlock *parser.RunConBlock, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission, l *logger.Logger, errors *[]string) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, stmt := range conBlock.Statements {
		wg.Add(1)
		go func(s interface{}) {
			defer wg.Done()
			localErrors := []string{}
			RunStatement(s, globalData, globalPermissions, l, &localErrors)
			if len(localErrors) > 0 {
				mu.Lock()
				*errors = append(*errors, localErrors...)
				mu.Unlock()
			}
		}(stmt)
	}
	wg.Wait()
}


// RunTask executes a task and handles any errors
func RunTask(t *parser.Task, globalData map[string]*parser.Data, globalPermissions map[string]*parser.Permission, l *logger.Logger) error {
	// Simulate task execution time
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Execute the task using the executor package
	err := executor.ExecuteTask(t.AgentName, t, globalData, globalPermissions, l)
	if err != nil {
		return err
	}

	// Optionally print task information
	PrintTask(t)
	return nil
}

// PrintTask prints the task information to the console
func PrintTask(t *parser.Task) {
	fmt.Printf("Task: %s, Agent: %s, Parameters: %v\n", t.TaskName, t.AgentName, t.Parameters)
}
