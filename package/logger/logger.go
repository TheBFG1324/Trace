package logger

import (
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
	newLog := Log{
		timestamp: time.Now(),
		information: information,
	}
	l.Logs = append(l.Logs, newLog)
}

func (l *Logger) GetLog(index int) Log {
	if index >=0 && index < len(l.Logs) {
		return l.Logs[index]
	}
	return Log{}
}

func (l *Logger) GetAllLogs() []Log {
	return l.Logs
}