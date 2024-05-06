package logger

// HttpMessageLoggerConfiguration represents options for logging HTTP message details.
type HttpMessageLoggerConfiguration struct {
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

// HttpMessageLoggerOptions represents a function type that can be used to apply configuration to the HttpMessageLoggerOptions struct.
type HttpMessageLoggerOptions func(*HttpMessageLoggerConfiguration)

// Default logger configuration
func defaultHttpMessageLoggerConfiguration() HttpMessageLoggerConfiguration {
	return HttpMessageLoggerConfiguration{
		logBody:            false,
		logHeaders:         false,
		headersToExclude:   []string{},
		headersToInclude:   []string{},
		headersToWhitelist: []string{},
	}
}

// NewHttpMessageLoggerConfiguration creates default HttpMessageLoggerConfiguration with the provided options.
func NewHttpMessageLoggerConfiguration(options ...HttpMessageLoggerOptions) HttpMessageLoggerConfiguration {
	config := defaultHttpMessageLoggerConfiguration()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithLogBody is an option that sets that enable to log body in the LoggingOptions.
func WithLogBody() HttpMessageLoggerOptions {
	return func(l *HttpMessageLoggerConfiguration) {
		l.logBody = true
	}
}

// WithLogHeaders is an option that sets that enable to log headers in the LoggingOptions.
func WithLogHeaders() HttpMessageLoggerOptions {
	return func(l *HttpMessageLoggerConfiguration) {
		l.logHeaders = true
	}
}

// WithHeadersToExclude is an option that sets the Headers To Exclude in the LoggingOptions.
func WithHeadersToExclude(headersToExclude ...string) HttpMessageLoggerOptions {
	return func(l *HttpMessageLoggerConfiguration) {
		l.headersToExclude = headersToExclude
	}
}

// WithHeadersToInclude is an option that sets the Headers To Include in the LoggingOptions.
func WithHeadersToInclude(headersToInclude ...string) HttpMessageLoggerOptions {
	return func(l *HttpMessageLoggerConfiguration) {
		l.headersToInclude = headersToInclude
	}
}

// WithHeadersToWhitelist is an option that sets the Headers To Whitelist in the LoggingOptions.
func WithHeadersToWhitelist(headersToWhitelist ...string) HttpMessageLoggerOptions {
	return func(l *HttpMessageLoggerConfiguration) {
		l.headersToWhitelist = headersToWhitelist
	}
}
