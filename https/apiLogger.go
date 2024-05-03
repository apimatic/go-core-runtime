package https

import (
	"net/http"
	"strings"
)

// ApiLoggerInterface Represents an interface for logging API requests and responses.
type ApiLoggerInterface interface {
	// Logs the details of an HTTP request.
	LogRequest(request *http.Request)
	// Logs the details of an HTTP response.
	LogResponse(response *http.Response)
}

// ApiLogger represents implementation for ApiLoggerInterface, providing methods to log HTTP requests and responses.
type ApiLogger struct {
	_loggingOptions LoggingOptions
	_logger         LoggerInterface
}

// NewApiLogger Constructs a new instance of ApiLogger.
func NewApiLogger(loggingOpt LoggingOptions) ApiLogger {
	return ApiLogger{
		_loggingOptions: loggingOpt,
		_logger:         loggingOpt.logger,
	}
}

// logRequest Logs an HTTP request.
func (a *ApiLogger) LogRequest(request *http.Request) {
	var logLevel = a._loggingOptions.logLevel
	var contentTypeHeader = a._getContentType(request.Header)
	var url string
	if a._loggingOptions.logRequest.includeQueryInPath {
		url = request.RequestURI
	} else {
		url = a._removeQueryParams(request.RequestURI)
	}

	a._logger.log(logLevel, "Request ${method} ${url} ${contentType}", map[string]any{
		"method":      request.Method,
		"url":         url,
		"contentType": contentTypeHeader,
	})

	a._applyLogRequestOptions(logLevel, request)
}

/**
 * Logs an HTTP response.
 * @param response The HTTP response to log.
 */
func (a *ApiLogger) LogResponse(response *http.Response) {
	var logLevel = a._loggingOptions.logLevel
	var contentTypeHeader = a._getContentType(response.Header)
	var contentLengthHeader = a._getContentLength(response.Header)

	a._logger.log(
		logLevel,
		"Response ${statusCode} ${contentLength} ${contentType}",
		map[string]any{
			"statusCode":    response.StatusCode,
			"contentLength": contentLengthHeader,
			"contentType":   contentTypeHeader,
		},
	)

	a._applyLogResponseOptions(logLevel, response)
}

func (a *ApiLogger) _applyLogRequestOptions(level LogLevel, request *http.Request) {
	a._applyLogRequestHeaders(
		level,
		request,
		a._loggingOptions.logRequest,
	)

	a._applyLogRequestBody(level, request, a._loggingOptions.logRequest)
}

func (a *ApiLogger) _applyLogRequestHeaders(
	level LogLevel,
	request *http.Request,
	logRequest HttpRequestLoggingOptions) {

	logHeaders := logRequest.logHeaders
	headersToInclude := logRequest.headersToInclude
	headersToExclude := logRequest.headersToExclude
	headersToWhitelist := logRequest.headersToWhitelist

	if logHeaders {
		var headersToLog = a._extractHeadersToLog(
			headersToInclude,
			headersToExclude,
			headersToWhitelist,
			request.Header,
		)

		a._logger.log(level, "Request headers ${headers}",
			map[string]any{"headers": headersToLog},
		)
	}
}

func (a *ApiLogger) _applyLogRequestBody(
	level LogLevel,
	request *http.Request,
	logRequest HttpRequestLoggingOptions) {

	if logRequest.logBody {
		a._logger.log(level, "Request body ${body}",
			map[string]any{"body": request.Body},
		)
	}
}

func (a *ApiLogger) _applyLogResponseOptions(level LogLevel, response *http.Response) {
	a._applyLogResponseHeaders(
		level,
		response,
		a._loggingOptions.logResponse,
	)

	a._applyLogResponseBody(
		level,
		response,
		a._loggingOptions.logResponse,
	)
}

func (a *ApiLogger) _applyLogResponseHeaders(
	level LogLevel,
	response *http.Response,
	logResponse HttpMessageLoggingOptions) {

	logHeaders := logResponse.logHeaders
	headersToInclude := logResponse.headersToInclude
	headersToExclude := logResponse.headersToExclude
	headersToWhitelist := logResponse.headersToWhitelist

	if logHeaders {
		var headersToLog = a._extractHeadersToLog(
			headersToInclude,
			headersToExclude,
			headersToWhitelist,
			response.Header,
		)

		a._logger.log(level, "Response headers ${headers}",
			map[string]any{"headers": headersToLog},
		)
	}
}

func (a *ApiLogger) _applyLogResponseBody(
	level LogLevel,
	response *http.Response,
	logResponse HttpMessageLoggingOptions) {

	if logResponse.logBody {
		a._logger.log(level, "Response body ${body}",
			map[string]any{"body": response.Body},
		)
	}
}

func (a *ApiLogger) _getContentType(headers http.Header) string {
	var contentType string = ""
	if len(headers) > 0 {
		contentType = headers.Get(CONTENT_TYPE_HEADER)
	}
	return contentType
}

func (a *ApiLogger) _getContentLength(headers http.Header) string {
	var contentLength string = ""
	if len(headers) > 0 {
		contentLength = headers.Get(CONTENT_LENGTH_HEADER)
	}
	return contentLength
}

func (a *ApiLogger) _removeQueryParams(url string) string {
	if strIndex := strings.Index(url, "?"); strIndex != -1 {
		return url[:strIndex]
	}
	return url
}

func (a *ApiLogger) _extractHeadersToLog(
	headersToInclude []string,
	headersToExclude []string,
	headersToWhitelist []string,
	headers http.Header) http.Header {
	filteredHeaders := http.Header{}
	if !(len(headers) > 0) {
		return http.Header{}
	}

	if len(headersToInclude) > 0 {
		filteredHeaders = a._includeHeadersToLog(
			headers,
			filteredHeaders,
			headersToInclude,
		)
	} else if len(headersToExclude) > 0 {
		filteredHeaders = a._excludeHeadersToLog(
			headers,
			filteredHeaders,
			headersToExclude,
		)
	} else {
		filteredHeaders = headers
	}

	return a._maskSenstiveHeaders(filteredHeaders, headersToWhitelist)
}

func (a *ApiLogger) _includeHeadersToLog(
	headers, filteredHeaders http.Header,
	headersToInclude []string) http.Header {
	// Filter headers based on the keys specified in headersToInclude
	for _, name := range headersToInclude {
		val, ok := headers[name]
		if len(val) > 0 && ok {
			filteredHeaders[name] = val
		}
	}
	return filteredHeaders
}

func (a *ApiLogger) _excludeHeadersToLog(
	headers, filteredHeaders http.Header,
	headersToExclude []string) http.Header {
	// Filter headers based on the keys specified in headersToExclude
	for key, value := range headers {
		if !_contains(key, headersToExclude) {
			if value != nil {
				filteredHeaders[key] = value
			}
		}
	}
	return filteredHeaders
}

func _contains(key string, slice []string) bool {
	for _, name := range slice {
		if name == key {
			return true
		}
	}
	return false
}

func (a *ApiLogger) _maskSenstiveHeaders(
	headers http.Header,
	headersToWhitelist []string) http.Header {

	if a._loggingOptions.maskSensitiveHeaders {
		for key, _ := range headers {
			val := headers.Get(key)
			headers.Set(key, a._maskIfSenstiveHeader(key, val, headersToWhitelist))
		}
	}
	return headers
}

func (a *ApiLogger) _maskIfSenstiveHeader(
	name, value string,
	headersToWhiteList []string) string {
	var nonSensitiveHeaders []string = []string{
		"accept",
		"accept-charset",
		"accept-encoding",
		"accept-language",
		"access-control-allow-origin",
		"cache-control",
		"connection",
		"content-encoding",
		"content-language",
		"content-length",
		"content-location",
		"content-md5",
		"content-range",
		"content-type",
		"date",
		"etag",
		"expect",
		"expires",
		"from",
		"host",
		"if-match",
		"if-modified-sinc",
		"if-none-match",
		"if-range",
		"if-unmodified-since",
		"keep-alive",
		"last-modified",
		"location",
		"max-forwards",
		"pragma",
		"range",
		"referer",
		"retry-after",
		"server",
		"trailer",
		"transfer-encoding",
		"upgrade",
		"user-agent",
		"vary",
		"via",
		"warning",
		"x-forwarded-for",
		"x-requested-with",
		"x-powered-by",
	}

	lowerCaseHeadersToWhiteList := make([]string, len(headersToWhiteList))
	for i, header := range headersToWhiteList {
		lowerCaseHeadersToWhiteList[i] = strings.ToLower(header)
	}

	if _contains(strings.ToLower(name), nonSensitiveHeaders) ||
		_contains(strings.ToLower(name), lowerCaseHeadersToWhiteList) {
		return value
	} else {
		return "**Redacted**"
	}
}
