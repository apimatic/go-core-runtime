package logger

// LoggerConfiguration represents options for configuring logging behavior.
type LoggerConfiguration struct {
	// The logger to use for logging messages.
	logger LoggerInterface
	// The log level to determine which messages should be logged.
	logLevel LogLevel
	// Options for logging HTTP requests.
	logRequest HttpRequestLoggerConfiguration
	// Options for logging HTTP responses.
	logResponse HttpMessageLoggerConfiguration
	// Indicates whether sensitive headers should be masked in logged messages.
	maskSensitiveHeaders bool
}

// LoggerOptions represents a function type that can be used to apply configuration to the LoggerOptions struct.
type LoggerOptions func(*LoggerConfiguration)

// Default logger configuration
func defaultLoggerConfiguration() LoggerConfiguration {
	return LoggerConfiguration{
		logger:               NullLogger{},
		logLevel:             LogLevel_INFO,
		logRequest:           NewHttpRequestLoggerConfiguration(),
		logResponse:          NewHttpMessageLoggerConfiguration(),
		maskSensitiveHeaders: true,
	}
}

// NewLoggerConfiguration creates default LoggingOptions with the provided options.
func NewLoggerConfiguration(options ...LoggerOptions) LoggerConfiguration {
	config := defaultLoggerConfiguration()

	for _, option := range options {
		option(&config)
	}
	return config
}

// WithLogger is an option that sets the LoggerInterface in the LoggingOptions.
func WithLogger(logger LoggerInterface) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.logger = logger
	}
}

// WithLogLevel is an option that sets the LogLevel in the LoggingOptions.
func WithLogLevel(level LogLevel) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.logLevel = level
	}
}

// WithMaskSensitiveHeaders is an option that enable to mask Sensitive Headers in the LoggingOptions.
func WithMaskSensitiveHeaders() LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.maskSensitiveHeaders = true
	}
}

// WithRequestConfiguration is an option that sets that enable to log Request in the LoggingOptions.
func WithRequestConfiguration(options ...HttpRequestLoggerOptions) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.logRequest = NewHttpRequestLoggerConfiguration(options...)
	}
}

// WithResponseConfiguration is an option that sets that enable to log Response in the LoggingOptions.
func WithResponseConfiguration(options ...HttpMessageLoggerOptions) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.logResponse = NewHttpMessageLoggerConfiguration(options...)
	}
}
