package https

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Represents options for configuring logging behavior.
type LoggingOptions struct {
	// The logger to use for logging messages.
	logger LoggerInterface
	// The log level to determine which messages should be logged.
	logLevel LogLevel
	// Options for logging HTTP requests.
	logRequest HttpRequestLoggingOptions
	// Options for logging HTTP responses.
	logResponse HttpMessageLoggingOptions
	// Indicates whether sensitive headers should be masked in logged messages.
	maskSensitiveHeaders bool
}

// Represents options for logging HTTP message details.
type HttpMessageLoggingOptions struct {
	// Indicates whether the message body should be logged.
	logBody bool
	// Indicates whether the message headers should be logged.
	logHeaders bool
	// Array of headers not to be displayed in logging.
	headersToExclude []string
	// Array of headers to be displayed in logging.
	headersToInclude []string
	// Array of headers which values are non-senstive to display in logging.
	headersToWhitelist []string
}

// Represents options for logging HTTP request details.
type HttpRequestLoggingOptions struct {
	HttpMessageLoggingOptions
	// Indicates whether the request query parameters should be included in the logged URL.
	includeQueryInPath bool
}

// LogLevel is a string enum.
// An enum representing different log levels.
type LogLevel string

// MarshalJSON implements the json.Marshaler interface for LogLevel.
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
// It customizes the JSON unmarshaling process for LogLevel objects.
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

// Default logging options
var DEFAULT_LOGGING_OPTIONS LoggingOptions = LoggingOptions{
	logger:   ConsoleLogger{},
	logLevel: LogLevel_INFO,
	logRequest: HttpRequestLoggingOptions{
		includeQueryInPath: false,
		HttpMessageLoggingOptions: HttpMessageLoggingOptions{
			logBody:            false,
			logHeaders:         false,
			headersToExclude:   []string{},
			headersToInclude:   []string{},
			headersToWhitelist: []string{},
		},
	},
	logResponse: HttpMessageLoggingOptions{
		logBody:            false,
		logHeaders:         false,
		headersToExclude:   []string{},
		headersToInclude:   []string{},
		headersToWhitelist: []string{},
	},
	maskSensitiveHeaders: true,
}

/*
// Create a new logging options object, falling back to the default values when not provided.
func mergeLoggingOptions(
	newOptions PartialLoggingOptions,
	defaultOptions LoggingOptions) LoggingOptions {
	//	return defaultsDeep({}, newOptions, defaultOptions);
	return defaultOptions
}
*/
