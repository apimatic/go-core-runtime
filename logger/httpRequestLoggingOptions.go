package logger

// HttpRequestLoggerConfiguration represents options for logging HTTP request details.
type HttpRequestLoggerConfiguration struct {
	HttpMessageLoggerConfiguration
	// Indicates whether the request query parameters should be included in the logged URL.
	includeQueryInPath bool
}

// HttpRequestLoggerOptions represents a function type that can be used to apply configuration to the HttpMessageLoggerOptions struct.
type HttpRequestLoggerOptions func(*HttpRequestLoggerConfiguration)

// Default logger configuration
func defaultHttpRequestLoggerConfiguration() HttpRequestLoggerConfiguration {
	return HttpRequestLoggerConfiguration{
		includeQueryInPath:             false,
		HttpMessageLoggerConfiguration: NewHttpMessageLoggerConfiguration(),
	}
}

// NewHttpRequestLoggerConfiguration creates default HttpRequestLoggerConfiguration with the provided options.
func NewHttpRequestLoggerConfiguration(options ...HttpRequestLoggerOptions) HttpRequestLoggerConfiguration {
	config := defaultHttpRequestLoggerConfiguration()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithIncludeQueryInPath is an option that enable include Query InPath in the LoggingOptions.
func WithIncludeQueryInPath() HttpRequestLoggerOptions {
	return func(l *HttpRequestLoggerConfiguration) {
		l.includeQueryInPath = true
	}
}
