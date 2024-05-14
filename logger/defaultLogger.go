package logger

import (
	"fmt"
	"regexp"
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

/**
 * Formats a message string by replacing placeholders with corresponding values from an object.
 * @param msg The message string containing placeholders.
 * @param obj The object containing values to replace placeholders.
 * @returns The formatted message string.
 */
func _formatMessage(msg string, obj map[string]any) string {
	// Use a regular expression to match placeholders in the message string
	regex := regexp.MustCompile(`\${([^}]+)}`)

	formattedMsg := regex.ReplaceAllStringFunc(msg, func(match string) string {
		key := match[2 : len(match)-1]
		if value, ok := obj[key]; ok {
			switch v := value.(type) {
			case string:
				return v
			default:
				return fmt.Sprintf("%v", v)
			}
		}
		return match
	})
	return formattedMsg
}
