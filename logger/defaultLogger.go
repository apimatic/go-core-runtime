package logger

import (
	"encoding/json"
	"fmt"
)

// LoggerInterface represents an interface for a generic logger.
type LoggerInterface interface {
	// Logs a message with a specified Log level and additional parameters.
	Log(level Level, message string, params map[string]any)
}

// ConsoleLogger represents a logger implementation that logs messages to the console.
type ConsoleLogger struct{}

// Logs a message to the console with the specified log level.
func (c ConsoleLogger) Log(level Level, message string, params map[string]any) {
	fmt.Println(level, ": ", _formatMessage(message, params))
}

func _formatMessage(message string, params map[string]any) string {
	byt, _ := json.Marshal(params)
	return message + " " + string(byt)
}
