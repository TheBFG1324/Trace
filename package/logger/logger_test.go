package logger_test

import (
	"testing"
	"trace/package/logger"
)

// TestNewLogger checks that the Logger is correctly initialized with a starting log entry.
func TestNewLogger(t *testing.T) {
	logger := logger.NewLogger()

	if len(logger.Logs) != 1 {
		t.Errorf("Expected 1 log entry on initialization, got %d", len(logger.Logs))
	}

	if logger.Logs[0].Information() != "Initialized Logger" {
		t.Errorf("Expected first log information to be 'Initialized Logger', got '%s'", logger.Logs[0].Information())
	}
}

// TestAddLog checks that a new log entry can be added and that it appears correctly.
func TestAddLog(t *testing.T) {
	logger := logger.NewLogger()
	initialCount := len(logger.Logs)

	logger.AddLog("Test log entry")

	if len(logger.Logs) != initialCount+1 {
		t.Errorf("Expected %d log entries, got %d", initialCount+1, len(logger.Logs))
	}

	lastLog := logger.Logs[len(logger.Logs)-1]
	if lastLog.Information() != "Test log entry" {
		t.Errorf("Expected last log information to be 'Test log entry', got '%s'", lastLog.Information())
	}
}

// TestGetLog checks that individual logs can be retrieved by index.
func TestGetLog(t *testing.T) {
	logger := logger.NewLogger()
	logger.AddLog("Second log entry")

	firstLog := logger.GetLog(0)
	if firstLog.Information() != "Initialized Logger" {
		t.Errorf("Expected first log information to be 'Initialized Logger', got '%s'", firstLog.Information())
	}

	secondLog := logger.GetLog(1)
	if secondLog.Information() != "Second log entry" {
		t.Errorf("Expected second log information to be 'Second log entry', got '%s'", secondLog.Information())
	}

	outOfBoundsLog := logger.GetLog(10)
	if !outOfBoundsLog.Timestamp().IsZero() || outOfBoundsLog.Information() != "" {
		t.Errorf("Expected out-of-bounds log to be empty, got '%v'", outOfBoundsLog)
	}
}

// TestGetAllLogs checks that all logs can be retrieved at once.
func TestGetAllLogs(t *testing.T) {
	logger := logger.NewLogger()
	logger.AddLog("Additional log 1")
	logger.AddLog("Additional log 2")

	allLogs := logger.GetAllLogs()
	expectedCount := 3

	if len(allLogs) != expectedCount {
		t.Errorf("Expected %d logs, got %d", expectedCount, len(allLogs))
	}

	if allLogs[1].Information() != "Additional log 1" || allLogs[2].Information() != "Additional log 2" {
		t.Errorf("Expected logs to match added information, got '%s' and '%s'", allLogs[1].Information(), allLogs[2].Information())
	}
}
