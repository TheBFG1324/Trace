package scheduler

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	"trace/package/parser"
)

// RunParentRequest schedules and runs the AICL parent request script
func RunParentRequest(p *parser.ParentRequest) bool {
	errors := []string{}
	statements := p.Statements
	for _, stmt := range statements {
		RunStatement(stmt)
	}

	if len(errors) != 0 {
		fmt.Println("Errors occured during runtime: ", errors)
		return false
	}
	return true
}

// RunStatement handles the execution of a single statement
func RunStatement(stmt interface{}) {
	switch s := stmt.(type) {
	case *parser.Task:
		RunTask(s)
	case *parser.RunSeqBlock:
		RunSeqBlock(s)
	case *parser.RunConBlock:
		RunConBlock(s)
	default:
		fmt.Println("Unknown statement type")
	}
}

// RunSeqBlock runs the tasks sequentially
func RunSeqBlock(seqBlock *parser.RunSeqBlock) {
	for _, stmt := range seqBlock.Statements {
		RunStatement(stmt)
	}
}

// RunConBlock runs the tasks concurrently
func RunConBlock(conBlock *parser.RunConBlock) {
	var wg sync.WaitGroup
	for _, stmt := range conBlock.Statements {
		wg.Add(1)
		go func(s interface{}) {
			defer wg.Done()
			RunStatement(s)
		}(stmt)
	}
	wg.Wait()
}

// RunTask simulates the execution of a task
func RunTask(t *parser.Task) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	PrintTask(t)
}

// PrintTask prints the task information to the console
func PrintTask(t *parser.Task) {
	fmt.Printf("Task: %s, Agent: %s, Parameters: %v\n", t.TaskName, t.AgentName, t.Parameters)
}