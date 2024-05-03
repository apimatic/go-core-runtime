package https

import (
	"encoding/json"
	"errors"
	"fmt"
)

// LoggingOptions represents options for configuring logging behavior.
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

// HttpMessageLoggingOptions represents options for logging HTTP message details.
type HttpMessageLoggingOptions struct {
	// Indicates whether the message body should be logged.
	logBody bool
	// Indicates whether the message headers should be logged.
	logHeaders bool
	// Array of headers not to be displayed in logging.
	headersToExclude []string
	// Array of headers to be displayed in logging.
	headersToInclude []string
	// Array of headers which values are non-sensitive to display in logging.
	headersToWhitelist []string
}

// HttpRequestLoggingOptions represents options for logging HTTP request details.
type HttpRequestLoggingOptions struct {
	HttpMessageLoggingOptions
	// Indicates whether the request query parameters should be included in the logged URL.
	includeQueryInPath bool
}

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

// LoggingConfigOptions represents a function type that can be used to apply configuration to the LoggingOptions struct.
type LoggingConfigOptions func(*LoggingOptions)

// Default logging options
func defaultLoggingOptions() LoggingOptions {
	return LoggingOptions{
		logger:   NullLogger{},
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
}

// NewLoggingOptions creates default LoggingOptions with the provided options.
func NewLoggingOptions(options ...LoggingConfigOptions) LoggingOptions {
	config := defaultLoggingOptions()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithLogger is an option that sets the LoggerInterface in the LoggingOptions.
func WithLogger(logger LoggerInterface) LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logger = logger
	}
}

// WithLogLevel is an option that sets the LogLevel in the LoggingOptions.
func WithLogLevel(level LogLevel) LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logLevel = level
	}
}

// WithMaskSensitiveHeaders is an option that enable to mask Sensitive Headers in the LoggingOptions.
func WithMaskSensitiveHeaders() LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.maskSensitiveHeaders = true
	}
}

// WithIncludeQueryInPath is an option that enable include Query InPath in the LoggingOptions.
func WithIncludeQueryInPath() LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logRequest.includeQueryInPath = true
	}
}

// WithLogRequestBody is an option that sets that enable to log Request body in the LoggingOptions.
func WithLogRequestBody() LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logRequest.logBody = true
	}
}

// WithLogRequestHeaders is an option that sets that enable to log Request headers in the LoggingOptions.
func WithLogRequestHeaders() LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logRequest.logHeaders = true
	}
}

// WithLogResponseBody is an option that sets that enable to log Request body in the LoggingOptions.
func WithLogResponseBody() LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logResponse.logBody = true
	}
}

// WithLogResponseHeaders is an option that sets that enable to log Request headers in the LoggingOptions.
func WithLogResponseHeaders() LoggingConfigOptions {
	return func(l *LoggingOptions) {
		l.logResponse.logHeaders = true
	}
}

// WithHeadersToExclude is an option that sets the Headers To Exclude in the LoggingOptions.
func WithHeadersToExclude(headersToExclude []string) LoggingConfigOptions {
	return func(l *LoggingOptions) {
		if l.logRequest.logHeaders {
			l.logRequest.headersToExclude = headersToExclude
		} else if l.logResponse.logHeaders {
			l.logResponse.headersToExclude = headersToExclude

		}
	}
}

// WithHeadersToInclude is an option that sets the Headers To Include in the LoggingOptions.
func WithHeadersToInclude(headersToInclude []string) LoggingConfigOptions {
	return func(l *LoggingOptions) {
		if l.logRequest.logHeaders {
			l.logRequest.headersToInclude = headersToInclude
		} else if l.logResponse.logHeaders {
			l.logResponse.headersToInclude = headersToInclude

		}
	}
}

// WithHeadersToWhitelist is an option that sets the Headers To Whitelist in the LoggingOptions.
func WithHeadersToWhitelist(headersToWhitelist []string) LoggingConfigOptions {
	return func(l *LoggingOptions) {
		if l.logRequest.logHeaders {
			l.logRequest.headersToWhitelist = headersToWhitelist
		} else if l.logResponse.logHeaders {
			l.logResponse.headersToWhitelist = headersToWhitelist

		}
	}
}
