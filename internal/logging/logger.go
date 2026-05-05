package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger is a thread-safe logger that writes JSON entries to a file
type Logger struct {
	file  *os.File
	mutex sync.Mutex
}

// New creates a new logger instance and initializes the log file
// It creates the logs directory if it doesn't exist
func New(logPath string) (*Logger, error) {
	// Create logs directory if it doesn't exist
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Open or create the log file in append mode
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		file:  file,
		mutex: sync.Mutex{},
	}, nil
}

// LogEvent logs an event with the given level, type, and data
func (l *Logger) LogEvent(level LogLevel, eventType EventType, data map[string]string) error {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		EventType: eventType,
		Data:      data,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Write to file with newline
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file == nil {
		return fmt.Errorf("logger is closed")
	}

	_, err = l.file.Write(append(jsonData, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// LogMessage is a convenience method to log a message sent event
func (l *Logger) LogMessage(sender, content string) error {
	return l.LogEvent(LevelInfo, EventMessageSent, MessageSentData(sender, content))
}

// LogError is a convenience method to log an error event
func (l *Logger) LogError(errorMsg string) error {
	return l.LogEvent(LevelError, EventError, ErrorData(errorMsg))
}

// Close closes the logger and flushes all pending writes
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file == nil {
		return fmt.Errorf("logger is already closed")
	}

	err := l.file.Close()
	l.file = nil
	return err
}
