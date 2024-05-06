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

// NewLoggerConfiguration creates default LoggingOptions with the provided options.
func NewLoggerConfiguration() LoggerConfiguration {
	return LoggerConfiguration{
		logger:               NullLogger{},
		logLevel:             LogLevel_INFO,
		logRequest:           NewHttpRequestLoggerConfiguration(),
		logResponse:          NewHttpMessageLoggerConfiguration(),
		maskSensitiveHeaders: true,
	}
}

// WithLogger is an option that sets the LoggerInterface in the LoggingOptions.
func (l *LoggerConfiguration) WithLogger(logger LoggerInterface) *LoggerConfiguration {
		l.logger = logger
		return l
}

// WithLogLevel is an option that sets the LogLevel in the LoggingOptions.
func (l *LoggerConfiguration)  WithLogLevel(level LogLevel) *LoggerConfiguration {
	l.logLevel = level
	return l
}

// WithMaskSensitiveHeaders is an option that enable to mask Sensitive Headers in the LoggingOptions.
func (l *LoggerConfiguration)  WithMaskSensitiveHeaders(maskSensitiveHeaders bool) *LoggerConfiguration {
	l.maskSensitiveHeaders = maskSensitiveHeaders
	return l
}

// WithRequestConfiguration is an option that sets that enable to log Request in the LoggingOptions.
func (l *LoggerConfiguration)  WithRequestConfiguration() *LoggerConfiguration {
	l.logRequest = NewHttpRequestLoggerConfiguration()
	return l
}

// WithResponseConfiguration is an option that sets that enable to log Response in the LoggingOptions.
func (l *LoggerConfiguration)  WithResponseConfiguration() *LoggerConfiguration {
	l.logResponse = NewHttpMessageLoggerConfiguration()
	return l
}
