package logging

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLoggerCreatesLogsDirectory(t *testing.T) {
	// Clean up before test
	defer os.RemoveAll("logs")

	logger, err := New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Verify logs directory exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		t.Error("logs directory was not created")
	}
}

func TestLoggerWritesJSONToFile(t *testing.T) {
	// Clean up before test
	defer os.RemoveAll("logs")

	logger, err := New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Log an event
	data := ClientJoinedData("TestUser", "127.0.0.1", "8080")
	err = logger.LogEvent(LevelInfo, EventClientJoined, data)
	if err != nil {
		t.Fatalf("Failed to log event: %v", err)
	}

	// Read the file and verify JSON
	content, err := os.ReadFile("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	var entry LogEntry
	err = json.Unmarshal(content, &entry)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if entry.Level != LevelInfo {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
	if entry.EventType != EventClientJoined {
		t.Errorf("Expected event ClientJoined, got %s", entry.EventType)
	}
	if entry.Data["clientName"] != "TestUser" {
		t.Errorf("Expected clientName TestUser, got %s", entry.Data["clientName"])
	}
}

func TestLoggerAppendsToExistingFile(t *testing.T) {
	// Clean up before test
	defer os.RemoveAll("logs")

	logger, err := New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Log first event
	data1 := ServerStartedData("8080")
	err = logger.LogEvent(LevelInfo, EventServerStarted, data1)
	if err != nil {
		t.Fatalf("Failed to log first event: %v", err)
	}

	logger.Close()

	// Create new logger instance (simulating server restart)
	logger, err = New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to recreate logger: %v", err)
	}
	defer logger.Close()

	// Log second event
	data2 := ClientJoinedData("User1", "127.0.0.1", "9000")
	err = logger.LogEvent(LevelInfo, EventClientJoined, data2)
	if err != nil {
		t.Fatalf("Failed to log second event: %v", err)
	}

	// Read the file and verify both entries
	content, err := os.ReadFile("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := parseJSONLines(string(content))
	if len(lines) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(lines))
	}

	// Verify first entry
	if lines[0].EventType != EventServerStarted {
		t.Errorf("Expected first event to be ServerStarted, got %s", lines[0].EventType)
	}

	// Verify second entry
	if lines[1].EventType != EventClientJoined {
		t.Errorf("Expected second event to be ClientJoined, got %s", lines[1].EventType)
	}
}

func TestLoggerHandlesConcurrentWrites(t *testing.T) {
	// Clean up before test
	defer os.RemoveAll("logs")

	logger, err := New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write events concurrently from multiple goroutines
	var wg sync.WaitGroup
	numGoroutines := 10
	eventsPerGoroutine := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				data := MessageSentData("User1", "Hello from goroutine")
				err := logger.LogEvent(LevelInfo, EventMessageSent, data)
				if err != nil {
					t.Errorf("Failed to log event: %v", err)
				}
				// Small delay to interleave writes
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// Read the file and verify all entries are valid JSON
	content, err := os.ReadFile("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := parseJSONLines(string(content))
	expectedCount := numGoroutines * eventsPerGoroutine
	if len(lines) != expectedCount {
		t.Errorf("Expected %d log entries, got %d", expectedCount, len(lines))
	}

	// Verify all entries are valid
	for i, entry := range lines {
		if entry.EventType != EventMessageSent {
			t.Errorf("Entry %d: Expected MessageSent, got %s", i, entry.EventType)
		}
		if entry.Data["sender"] != "User1" {
			t.Errorf("Entry %d: Expected sender User1, got %s", i, entry.Data["sender"])
		}
	}
}

func TestLoggerTimestampFormat(t *testing.T) {
	// Clean up before test
	defer os.RemoveAll("logs")

	logger, err := New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	data := ServerStartedData("8080")
	err = logger.LogEvent(LevelInfo, EventServerStarted, data)
	if err != nil {
		t.Fatalf("Failed to log event: %v", err)
	}

	content, err := os.ReadFile("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	var entry LogEntry
	err = json.Unmarshal(content, &entry)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify timestamp is in RFC3339 format and parseable
	_, err = time.Parse(time.RFC3339, entry.Timestamp)
	if err != nil {
		t.Errorf("Timestamp is not in RFC3339 format: %s", entry.Timestamp)
	}
}

func TestLoggerClose(t *testing.T) {
	// Clean up before test
	defer os.RemoveAll("logs")

	logger, err := New("logs/test.log")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	data := ServerStartedData("8080")
	err = logger.LogEvent(LevelInfo, EventServerStarted, data)
	if err != nil {
		t.Fatalf("Failed to log event: %v", err)
	}

	// Close should not error
	err = logger.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Logging after close should error
	err = logger.LogEvent(LevelInfo, EventServerStarted, data)
	if err == nil {
		t.Error("Expected error when logging after close, got nil")
	}
}

// Helper function to parse JSON lines from file content
func parseJSONLines(content string) []LogEntry {
	var entries []LogEntry
	var currentLine string

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			if currentLine != "" {
				var entry LogEntry
				if err := json.Unmarshal([]byte(currentLine), &entry); err == nil {
					entries = append(entries, entry)
				}
				currentLine = ""
			}
		} else {
			currentLine += string(content[i])
		}
	}

	// Don't forget the last line if there's no trailing newline
	if currentLine != "" {
		var entry LogEntry
		if err := json.Unmarshal([]byte(currentLine), &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	return entries
}
