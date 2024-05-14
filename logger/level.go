package logger

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Level is a string enum.
// An enum representing different log levels.
type Level string

// MarshalJSON implements the json.Marshaller interface for Level.
// It customizes the JSON marshaling process for Level objects.
func (e Level) MarshalJSON() (
	[]byte,
	error) {
	if e.isValid() {
		return []byte(fmt.Sprintf("\"%v\"", e)), nil
	}
	return nil, errors.New("the provided enum value is not allowed for Level")
}

// UnmarshalJSON implements the json.Unmarshaler interface for Level.
// It customizes the JSON unmarshalling process for Level objects.
func (e *Level) UnmarshalJSON(input []byte) error {
	var enumValue string
	err := json.Unmarshal(input, &enumValue)
	if err != nil {
		return err
	}
	*e = Level(enumValue)
	if !e.isValid() {
		return errors.New("the value " + string(input) + " cannot be unmarshalled to Level")
	}
	return nil
}

// Checks whether the value is actually a member of Level.
func (e Level) isValid() bool {
	switch e {
	case Level_ERROR,
		Level_WARN,
		Level_INFO,
		Level_DEBUG,
		Level_TRACE:
		return true
	}
	return false
}

const (
	Level_ERROR Level = "error" // Error log level.
	Level_WARN  Level = "warn"  // Warning log level.
	Level_INFO  Level = "info"  // Information log level.
	Level_DEBUG Level = "debug" // Debug log level.
	Level_TRACE Level = "trace" // Trace log level.
)
