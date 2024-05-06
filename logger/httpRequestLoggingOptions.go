package logger

// LoggerRequestConfiguration represents options for logging HTTP request details.
type LoggerRequestConfiguration struct {
	LoggerResponseConfiguration
	// Indicates whether the request query parameters should be included in the logged URL.
	includeQueryInPath bool
}

// NewLoggerRequestConfiguration creates default HttpRequestLoggerConfiguration with the provided options.
func NewLoggerRequestConfiguration() *LoggerRequestConfiguration {
	return &LoggerRequestConfiguration{
		includeQueryInPath:             false,
		LoggerResponseConfiguration: *NewLoggerResponseConfiguration(),
	}
}

// WithIncludeQueryInPath is an option that enable include Query InPath in the LoggingOptions.
func (h *LoggerRequestConfiguration) WithIncludeQueryInPath(includeQueryInPath bool) *LoggerRequestConfiguration {
	h.includeQueryInPath = includeQueryInPath
	return h
}


// WithLogBody is an option that sets that enable to log body in the LoggingOptions.
func (h *LoggerRequestConfiguration) WithLogBody(logBody bool) *LoggerRequestConfiguration {
	h.logBody = logBody
	return h
}

// WithLogHeaders is an option that sets that enable to log headers in the LoggingOptions.
func (h *LoggerRequestConfiguration) WithLogHeaders(logHeaders bool) *LoggerRequestConfiguration {
	h.logHeaders = logHeaders
	return h
}

// WithHeadersToExclude is an option that sets the Headers To Exclude in the LoggingOptions.
func (h *LoggerRequestConfiguration) WithHeadersToExclude(headersToExclude ...string) *LoggerRequestConfiguration {
	h.headersToExclude = headersToExclude
	return h
}

// WithHeadersToInclude is an option that sets the Headers To Include in the LoggingOptions.
func (h *LoggerRequestConfiguration) WithHeadersToInclude(headersToInclude ...string) *LoggerRequestConfiguration {
	h.headersToInclude = headersToInclude
	return h
}

// WithHeadersToWhitelist is an option that sets the Headers To Whitelist in the LoggingOptions.
func (h *LoggerRequestConfiguration) WithHeadersToWhitelist(headersToWhitelist ...string) *LoggerRequestConfiguration {
	h.headersToWhitelist = headersToWhitelist
	return h
}