package logger

// ResponseLoggerConfiguration represents options for logging HTTP message details.
type ResponseLoggerConfiguration struct {
	messageLoggerConfiguration
}

// ResponseLoggerOptions represents a function type that can be used to apply configuration to the ResponseLoggerOptions struct.
type ResponseLoggerOptions func(*ResponseLoggerConfiguration)

// Default logger configuration
func defaultResponseLoggerConfiguration() ResponseLoggerConfiguration {
	return ResponseLoggerConfiguration{
		messageLoggerConfiguration: defaultMessageLoggerConfiguration(),
	}
}

// NewResponseLoggerConfiguration creates default ResponseLoggerConfiguration with the provided options.
func NewResponseLoggerConfiguration(options ...ResponseLoggerOptions) ResponseLoggerConfiguration {
	config := defaultResponseLoggerConfiguration()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithResponseBody is an option that sets that enable to log body in the LoggingOptions.
func WithResponseBody(logBody bool) ResponseLoggerOptions {
	return func(l *ResponseLoggerConfiguration) {
		l.body = logBody
	}
}

// WithResponseHeaders is an option that sets that enable to log headers in the LoggingOptions.
func WithResponseHeaders(logHeaders bool) ResponseLoggerOptions {
	return func(l *ResponseLoggerConfiguration) {
		l.headers = logHeaders
	}
}

// WithExcludeResponseHeaders is an option that sets the Headers To Exclude in the LoggingOptions.
func WithExcludeResponseHeaders(excludeHeaders ...string) ResponseLoggerOptions {
	return func(l *ResponseLoggerConfiguration) {
		l.excludeHeaders = excludeHeaders
	}
}

// WithIncludeResponseHeaders is an option that sets the Headers To Include in the LoggingOptions.
func WithIncludeResponseHeaders(includeHeaders ...string) ResponseLoggerOptions {
	return func(l *ResponseLoggerConfiguration) {
		l.includeHeaders = includeHeaders
	}
}

// WithWhitelistResponseHeaders is an option that sets the Headers To Whitelist in the LoggingOptions.
func WithWhitelistResponseHeaders(whitelistHeaders ...string) ResponseLoggerOptions {
	return func(l *ResponseLoggerConfiguration) {
		l.whitelistHeaders = whitelistHeaders
	}
}
