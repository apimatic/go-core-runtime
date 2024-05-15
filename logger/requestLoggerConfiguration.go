package logger

// RequestLoggerConfiguration represents options for logging HTTP request details.
type RequestLoggerConfiguration struct {
	MessageLoggerConfiguration
	// Indicates whether the request query parameters should be included in the logged URL.
	includeQueryInPath bool
}

// RequestLoggerOptions represents a function type that can be used to apply configuration to the RequestLoggerOptions struct.
type RequestLoggerOptions func(*RequestLoggerConfiguration)

// Default logger configuration
func defaultRequestLoggerConfiguration() RequestLoggerConfiguration {
	return RequestLoggerConfiguration{
		includeQueryInPath:         false,
		MessageLoggerConfiguration: defaultMessageLoggerConfiguration(),
	}
}

// NewHttpRequestLoggerConfiguration creates default RequestLoggerConfiguration with the provided options.
func NewHttpRequestLoggerConfiguration(options ...RequestLoggerOptions) RequestLoggerConfiguration {
	config := defaultRequestLoggerConfiguration()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithIncludeQueryInPath is an option that enable include Query InPath in the LoggingOptions.
func WithIncludeQueryInPath(includeQueryInPath bool) RequestLoggerOptions {
	return func(l *RequestLoggerConfiguration) {
		l.includeQueryInPath = includeQueryInPath
	}
}

// WithRequestBody is an option that sets that enable to log body in the LoggingOptions.
func WithRequestBody(logBody bool) RequestLoggerOptions {
	return func(l *RequestLoggerConfiguration) {
		l.body = logBody
	}
}

// WithRequestHeaders is an option that sets that enable to log headers in the LoggingOptions.
func WithRequestHeaders(logHeaders bool) RequestLoggerOptions {
	return func(l *RequestLoggerConfiguration) {
		l.headers = logHeaders
	}
}

// WithExcludeRequestHeaders is an option that sets the Headers To Exclude in the LoggingOptions.
func WithExcludeRequestHeaders(excludeHeaders ...string) RequestLoggerOptions {
	return func(l *RequestLoggerConfiguration) {
		l.excludeHeaders = excludeHeaders
	}
}

// WithIncludeRequestHeaders is an option that sets the Headers To Include in the LoggingOptions.
func WithIncludeRequestHeaders(includeHeaders ...string) RequestLoggerOptions {
	return func(l *RequestLoggerConfiguration) {
		l.includeHeaders = includeHeaders
	}
}

// WithWhitelistRequestHeaders is an option that sets the Headers To Whitelist in the LoggingOptions.
func WithWhitelistRequestHeaders(whitelistHeaders ...string) RequestLoggerOptions {
	return func(l *RequestLoggerConfiguration) {
		l.whitelistHeaders = whitelistHeaders
	}
}
