package logger

// HttpRequestLoggerConfiguration represents options for logging HTTP request details.
type HttpRequestLoggerConfiguration struct {
	HttpMessageLoggerConfiguration
	// Indicates whether the request query parameters should be included in the logged URL.
	includeQueryInPath bool
}

// NewHttpRequestLoggerConfiguration creates default HttpRequestLoggerConfiguration with the provided options.
func NewHttpRequestLoggerConfiguration() HttpRequestLoggerConfiguration {
	return HttpRequestLoggerConfiguration{
		includeQueryInPath:             false,
		HttpMessageLoggerConfiguration: NewHttpMessageLoggerConfiguration(),
	}
}

// WithIncludeQueryInPath is an option that enable include Query InPath in the LoggingOptions.
func (h *HttpRequestLoggerConfiguration) WithIncludeQueryInPath(includeQueryInPath bool) *HttpRequestLoggerConfiguration {
	h.includeQueryInPath = includeQueryInPath
	return h
}


// WithLogBody is an option that sets that enable to log body in the LoggingOptions.
func (h *HttpRequestLoggerConfiguration) WithLogBody(logBody bool) *HttpRequestLoggerConfiguration {
	h.logBody = logBody
	return h
}

// WithLogHeaders is an option that sets that enable to log headers in the LoggingOptions.
func (h *HttpRequestLoggerConfiguration) WithLogHeaders(logHeaders bool) *HttpRequestLoggerConfiguration {
	h.logHeaders = logHeaders
	return h
}

// WithHeadersToExclude is an option that sets the Headers To Exclude in the LoggingOptions.
func (h *HttpRequestLoggerConfiguration) WithHeadersToExclude(headersToExclude ...string) *HttpRequestLoggerConfiguration {
	h.headersToExclude = headersToExclude
	return h
}

// WithHeadersToInclude is an option that sets the Headers To Include in the LoggingOptions.
func (h *HttpRequestLoggerConfiguration) WithHeadersToInclude(headersToInclude ...string) *HttpRequestLoggerConfiguration {
	h.headersToInclude = headersToInclude
	return h
}

// WithHeadersToWhitelist is an option that sets the Headers To Whitelist in the LoggingOptions.
func (h *HttpRequestLoggerConfiguration) WithHeadersToWhitelist(headersToWhitelist ...string) *HttpRequestLoggerConfiguration {
	h.headersToWhitelist = headersToWhitelist
	return h
}