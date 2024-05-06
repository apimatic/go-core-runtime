package logger

// LoggerResponseConfiguration represents options for logging HTTP message details.
type LoggerResponseConfiguration struct {
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

// NewLoggerResponseConfiguration creates default HttpMessageLoggerConfiguration with the provided options.
func NewLoggerResponseConfiguration() *LoggerResponseConfiguration {
	return &LoggerResponseConfiguration{
		logBody:            false,
		logHeaders:         false,
		headersToExclude:   []string{},
		headersToInclude:   []string{},
		headersToWhitelist: []string{},
	}
}

// WithLogBody is an option that sets that enable to log body in the LoggingOptions.
func (h *LoggerResponseConfiguration) WithLogBody(logBody bool) *LoggerResponseConfiguration {
	h.logBody = logBody
	return h
}

// WithLogHeaders is an option that sets that enable to log headers in the LoggingOptions.
func (h *LoggerResponseConfiguration) WithLogHeaders(logHeaders bool) *LoggerResponseConfiguration {
	h.logHeaders = logHeaders
	return h
}

// WithHeadersToExclude is an option that sets the Headers To Exclude in the LoggingOptions.
func (h *LoggerResponseConfiguration) WithHeadersToExclude(headersToExclude ...string) *LoggerResponseConfiguration {
	h.headersToExclude = headersToExclude
	return h
}

// WithHeadersToInclude is an option that sets the Headers To Include in the LoggingOptions.
func (h *LoggerResponseConfiguration) WithHeadersToInclude(headersToInclude ...string) *LoggerResponseConfiguration {
	h.headersToInclude = headersToInclude
	return h
}

// WithHeadersToWhitelist is an option that sets the Headers To Whitelist in the LoggingOptions.
func (h *LoggerResponseConfiguration) WithHeadersToWhitelist(headersToWhitelist ...string) *LoggerResponseConfiguration {
	h.headersToWhitelist = headersToWhitelist
	return h
}
