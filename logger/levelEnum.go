package logger

import (
	"encoding/json"
	"errors"
	"fmt"
)

// LogLevel is a string enum.
// An enum representing different log levels.
type LogLevel string

// MarshalJSON implements the json.Marshaller interface for LogLevel.
// It customizes the JSON marshaling process for LogLevel objects.
func (e LogLevel) MarshalJSON() (
	[]byte,
	error) {
	if e.isValid() {
		return []byte(fmt.Sprintf("\"%v\"", e)), nil
	}
	return nil, errors.New("the provided enum value is not allowed for LogLevel")
}

// UnmarshalJSON implements the json.Unmarshaler interface for LogLevel.
// It customizes the JSON unmarshalling process for LogLevel objects.
func (e *LogLevel) UnmarshalJSON(input []byte) error {
	var enumValue string
	err := json.Unmarshal(input, &enumValue)
	if err != nil {
		return err
	}
	*e = LogLevel(enumValue)
	if !e.isValid() {
		return errors.New("the value " + string(input) + " cannot be unmarshalled to LogLevel")
	}
	return nil
}

// Checks whether the value is actually a member of LogLevel.
func (e LogLevel) isValid() bool {
	switch e {
	case LogLevel_ERROR,
		LogLevel_WARN,
		LogLevel_INFO,
		LogLevel_DEBUG,
		LogLevel_TRACE:
		return true
	}
	return false
}

const (
	LogLevel_ERROR LogLevel = "error" // Error log level.
	LogLevel_WARN  LogLevel = "warn"  // Warning log level.
	LogLevel_INFO  LogLevel = "info"  // Information log level.
	LogLevel_DEBUG LogLevel = "debug" // Debug log level.
	LogLevel_TRACE LogLevel = "trace" // Trace log level.
)
