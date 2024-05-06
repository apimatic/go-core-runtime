package logger

import (
	"fmt"
)

// LoggerInterface represents an interface for a generic logger.
type LoggerInterface interface {
	// Logs a message with a specified Log level and additional parameters.
	Log(level LogLevel, message string, params map[string]any)
}

// ConsoleLogger represents a logger implementation that logs messages to the console.
type ConsoleLogger struct{}

// Logs a message to the console with the specified log level.
func (c ConsoleLogger) Log(level LogLevel, message string, params map[string]any) {
	fmt.Println(level, ": ", message)
}

// NullLogger represents a logger implementation that does not perform any logging.
// Messages logged using this logger are effectively ignored.
type NullLogger struct{}

// Logs a message. Since this is a null logger, the log method does nothing.
func (n NullLogger) Log(level LogLevel, _message string, _params map[string]any) {
	// This is a null logger, so it does not perform any logging.
	// All parameters are ignored.
	return
}
