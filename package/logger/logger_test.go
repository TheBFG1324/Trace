package logger_test

import (
	"testing"
	"trace/package/logger"
)

// TestNewLogger checks that the Logger is correctly initialized with a starting log entry.
func TestNewLogger(t *testing.T) {
	log := logger.NewLogger()

	if len(log.Logs) != 1 {
		t.Errorf("Expected 1 log entry on initialization, got %d", len(log.Logs))
	}

	if log.Logs[0].Information() != "Initialized Logger" {
		t.Errorf("Expected first log information to be 'Initialized Logger', got '%s'", log.Logs[0].Information())
	}
}

// TestAddLog checks that a new log entry can be added and that it appears correctly.
func TestAddLog(t *testing.T) {
	log := logger.NewLogger()
	initialCount := len(log.Logs)

	newLog := logger.NewLog("Test log entry")
	log.AddLog(newLog)

	if len(log.Logs) != initialCount+1 {
		t.Errorf("Expected %d log entries, got %d", initialCount+1, len(log.Logs))
	}

	lastLog := log.Logs[len(log.Logs)-1]
	if lastLog.Information() != "Test log entry" {
		t.Errorf("Expected last log information to be 'Test log entry', got '%s'", lastLog.Information())
	}
}

// TestGetLog checks that individual logs can be retrieved by index.
func TestGetLog(t *testing.T) {
	log := logger.NewLogger()
	log.AddLog(logger.NewLog("Second log entry"))

	firstLog := log.GetLog(0)
	if firstLog.Information() != "Initialized Logger" {
		t.Errorf("Expected first log information to be 'Initialized Logger', got '%s'", firstLog.Information())
	}

	secondLog := log.GetLog(1)
	if secondLog.Information() != "Second log entry" {
		t.Errorf("Expected second log information to be 'Second log entry', got '%s'", secondLog.Information())
	}

	outOfBoundsLog := log.GetLog(10)
	if !outOfBoundsLog.Timestamp().IsZero() || outOfBoundsLog.Information() != "" {
		t.Errorf("Expected out-of-bounds log to be empty, got '%v'", outOfBoundsLog)
	}
}

// TestGetAllLogs checks that all logs can be retrieved at once.
func TestGetAllLogs(t *testing.T) {
	log := logger.NewLogger()
	log.AddLog(logger.NewLog("Additional log 1"))
	log.AddLog(logger.NewLog("Additional log 2"))

	allLogs := log.GetAllLogs()
	expectedCount := 3

	if len(allLogs) != expectedCount {
		t.Errorf("Expected %d logs, got %d", expectedCount, len(allLogs))
	}

	if allLogs[1].Information() != "Additional log 1" || allLogs[2].Information() != "Additional log 2" {
		t.Errorf("Expected logs to match added information, got '%s' and '%s'", allLogs[1].Information(), allLogs[2].Information())
	}
}
