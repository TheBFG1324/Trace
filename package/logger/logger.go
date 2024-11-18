package logger

import (
	"sync"
	"time"
)

type Log struct {
	timestamp	time.Time
	information	string
}

func (l *Log) Timestamp() time.Time {
	return l.timestamp
}

func (l *Log) Information() string {
	return l.information
}

type Logger struct {
	Logs	[]Log
	mu 		sync.Mutex
}

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

func(l *Logger) AddLog(information string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	newLog := Log{
		timestamp: time.Now(),
		information: information,
	}
	l.Logs = append(l.Logs, newLog)
}

func (l *Logger) GetLog(index int) Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index >=0 && index < len(l.Logs) {
		return l.Logs[index]
	}
	return Log{}
}

func (l *Logger) GetAllLogs() []Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	logsCopy := make([]Log, len(l.Logs))
	copy(logsCopy, l.Logs)
	return logsCopy
}