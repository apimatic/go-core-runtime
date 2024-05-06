package logger

// LoggerConfiguration represents options for configuring logging behavior.
type LoggerConfiguration struct {
	// The logger to use for logging messages.
	logger LoggerInterface
	// The log level to determine which messages should be logged.
	level Level
	// Options for logging HTTP requests.
	request RequestLoggerConfiguration
	// Options for logging HTTP responses.
	response MessageLoggerConfiguration
	// Indicates whether sensitive headers should be masked in logged messages.
	maskSensitiveHeaders bool
}

// LoggerOptions represents a function type that can be used to apply configuration to the LoggerOptions struct.
type LoggerOptions func(*LoggerConfiguration)

// Default logger configuration
func defaultLoggerConfiguration() LoggerConfiguration {
	return LoggerConfiguration{
		logger:               NullLogger{},
		level:                LogLevel_INFO,
		request:              NewHttpRequestLoggerConfiguration(),
		response:             NewResponseLoggerConfiguration(),
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

// WithLevel is an option that sets the LogLevel in the LoggingOptions.
func WithLevel(level Level) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.level = level
	}
}

// WithMaskSensitiveHeaders is an option that enable to mask Sensitive Headers in the LoggingOptions.
func WithMaskSensitiveHeaders() LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.maskSensitiveHeaders = true
	}
}

// WithRequestConfiguration is an option that sets that enable to log Request in the LoggingOptions.
func WithRequestConfiguration(options ...RequestLoggerOptions) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.request = NewHttpRequestLoggerConfiguration(options...)
	}
}

// WithResponseConfiguration is an option that sets that enable to log Response in the LoggingOptions.
func WithResponseConfiguration(options ...MessageLoggerOptions) LoggerOptions {
	return func(l *LoggerConfiguration) {
		l.response = NewResponseLoggerConfiguration(options...)
	}
}
