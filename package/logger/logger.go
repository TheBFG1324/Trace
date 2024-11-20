package logger

import (
	"sync"
	"time"
	"fmt"
)

// Log struct holds information and a timestamp
type Log struct {
	timestamp	time.Time
	information	string
}

// NewLog takes in an information string and outputs a new Log struct
func NewLog(information string) Log {
	newLog := Log{
		timestamp: time.Now(),
		information: information,
	}
	return newLog
}

// Timestamp gets the timestamp of a Log
func (l *Log) Timestamp() time.Time {
	return l.timestamp
}

// Information gets the information string of a Log
func (l *Log) Information() string {
	return l.information
}

type Logger struct {
	Logs	[]Log
	mu 		sync.Mutex
}

// NewLogger returns a pointer to a new logger struct
func NewLogger() *Logger {
	newLog := Log{
		timestamp: time.Now(),
		information: "Initialized Logger",
	}
	newLogger := &Logger{
		Logs: []Log{newLog},
	}
	return newLogger
}

// AddLog adds a Log to a given logger
func (l *Logger) AddLog(log Log) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logs = append(l.Logs, log)
}

// AddLogs adds multiple logs at once
func (l *Logger) AddLogs(logs []Log) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logs = append(l.Logs, logs...)
}

// GetLog gets the Log at a given idex value
func (l *Logger) GetLog(index int) Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index >=0 && index < len(l.Logs) {
		return l.Logs[index]
	}
	return Log{}
}

// GetAllLogs returns all of the logs stored in a logger
func (l *Logger) GetAllLogs() []Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	logsCopy := make([]Log, len(l.Logs))
	copy(logsCopy, l.Logs)
	return logsCopy
}

// PrintAllLogs prints all logs stored in the logger.
func (l *Logger) PrintAllLogs() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, log := range l.Logs {
		fmt.Printf("[%s] %s\n", log.timestamp.Format(time.RFC3339), log.information)
	}
}