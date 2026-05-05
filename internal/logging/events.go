package logging

// LogLevel represents the severity of a log entry
type LogLevel string

const (
	LevelInfo    LogLevel = "INFO"
	LevelWarning LogLevel = "WARNING"
	LevelError   LogLevel = "ERROR"
	LevelDebug   LogLevel = "DEBUG"
)

// EventType represents the type of event being logged
type EventType string

const (
	EventServerStarted      EventType = "ServerStarted"
	EventServerStopped      EventType = "ServerStopped"
	EventClientJoined       EventType = "ClientJoined"
	EventClientDisconnected EventType = "ClientDisconnected"
	EventMessageSent        EventType = "MessageSent"
	EventNameChanged        EventType = "NameChanged"
	EventError              EventType = "Error"
)

// LogEntry represents a single log entry in JSON format
type LogEntry struct {
	Timestamp string            `json:"timestamp"`
	Level     LogLevel          `json:"level"`
	EventType EventType         `json:"eventType"`
	Data      map[string]string `json:"data"`
}

// EventData helper functions to create consistent data maps

func ServerStartedData(port string) map[string]string {
	return map[string]string{
		"port": port,
	}
}

func ServerStoppedData() map[string]string {
	return map[string]string{
		"message": "Server stopped",
	}
}

func ClientJoinedData(clientName, ip, port string) map[string]string {
	return map[string]string{
		"clientName": clientName,
		"ip":         ip,
		"port":       port,
	}
}

func ClientDisconnectedData(clientName string) map[string]string {
	return map[string]string{
		"clientName": clientName,
	}
}

func MessageSentData(sender, content string) map[string]string {
	return map[string]string{
		"sender":  sender,
		"content": content,
	}
}

func NameChangedData(oldName, newName string) map[string]string {
	return map[string]string{
		"oldName": oldName,
		"newName": newName,
	}
}

func ErrorData(errorMsg string) map[string]string {
	return map[string]string{
		"error": errorMsg,
	}
}
