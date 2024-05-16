package logger

import (
	"fmt"
	"regexp"
)

// LoggerInterface represents an interface for a generic logger.
type LoggerInterface interface {
	// Log function provides a message string containing placeholders in the format '%{key}',
	// along with the log level and a map of parameters that can be replaced in the message.
	Log(level Level, message string, params map[string]any)
}

// ConsoleLogger represents a logger implementation that logs messages to the console.
type ConsoleLogger struct{}

// Log function provides a message string containing placeholders in the format '%{key}',
// along with the log level and a map of parameters that can be replaced in the message.
func (c ConsoleLogger) Log(level Level, message string, params map[string]any) {
	fmt.Println(level, ": ", _formatMessage(message, params))
}

func _formatMessage(msg string, obj map[string]interface{}) string {
	regex := regexp.MustCompile(`\%{([^}]+)}`)

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
