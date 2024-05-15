package logger

// MessageLoggerConfiguration represents options for logging HTTP message details.
type MessageLoggerConfiguration struct {
	// Indicates whether the message body should be logged.
	body bool
	// Indicates whether the message headers should be logged.
	headers bool
	// Array of headers not to be displayed in logging.
	excludeHeaders []string
	// Array of headers to be displayed in logging.
	includeHeaders []string
	// Array of headers which values are non-sensitive to display in logging.
	whitelistHeaders []string
}

// MessageLoggerOptions represents a function type that can be used to apply configuration to the MessageLoggerOptions struct.
type MessageLoggerOptions func(*MessageLoggerConfiguration)

// Default logger configuration
func defaultMessageLoggerConfiguration() MessageLoggerConfiguration {
	return MessageLoggerConfiguration{
		body:             false,
		headers:          false,
		excludeHeaders:   []string{},
		includeHeaders:   []string{},
		whitelistHeaders: []string{},
	}
}

// NewResponseLoggerConfiguration creates default MessageLoggerConfiguration with the provided options.
func NewResponseLoggerConfiguration(options ...MessageLoggerOptions) MessageLoggerConfiguration {
	config := defaultMessageLoggerConfiguration()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithResponseBody is an option that sets that enable to log body in the LoggingOptions.
func WithResponseBody(logBody bool) MessageLoggerOptions {
	return func(l *MessageLoggerConfiguration) {
		l.body = logBody
	}
}

// WithResponseHeaders is an option that sets that enable to log headers in the LoggingOptions.
func WithResponseHeaders(logHeaders bool) MessageLoggerOptions {
	return func(l *MessageLoggerConfiguration) {
		l.headers = logHeaders
	}
}

// WithExcludeResponseHeaders is an option that sets the Headers To Exclude in the LoggingOptions.
func WithExcludeResponseHeaders(excludeHeaders ...string) MessageLoggerOptions {
	return func(l *MessageLoggerConfiguration) {
		l.excludeHeaders = excludeHeaders
	}
}

// WithIncludeResponseHeaders is an option that sets the Headers To Include in the LoggingOptions.
func WithIncludeResponseHeaders(includeHeaders ...string) MessageLoggerOptions {
	return func(l *MessageLoggerConfiguration) {
		l.includeHeaders = includeHeaders
	}
}

// WithWhitelistResponseHeaders is an option that sets the Headers To Whitelist in the LoggingOptions.
func WithWhitelistResponseHeaders(whitelistHeaders ...string) MessageLoggerOptions {
	return func(l *MessageLoggerConfiguration) {
		l.whitelistHeaders = whitelistHeaders
	}
}
