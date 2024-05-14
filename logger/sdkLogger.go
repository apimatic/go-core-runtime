package logger

import (
	"fmt"
	"net/http"
	"strings"
)

// SdkLoggerInterface Represents an interface for logging API requests and responses.
type SdkLoggerInterface interface {
	// LogRequest logs the details of an HTTP request.
	LogRequest(request *http.Request)
	// LogResponse logs the details of an HTTP response.
	LogResponse(response *http.Response)
}

// SdkLogger represents implementation for SdkLoggerInterface, providing methods to log HTTP requests and responses.
type SdkLogger struct {
	loggingOptions LoggerConfiguration
	logger         LoggerInterface
}

// NewSdkLogger Constructs a new instance of SdkLogger.
func NewSdkLogger(loggingOpt LoggerConfiguration) *SdkLogger {
	return &SdkLogger{
		loggingOptions: loggingOpt,
		logger:         loggingOpt.logger,
	}
}

// LogRequest request Logs an HTTP request.
func (a *SdkLogger) LogRequest(request *http.Request) {
	var logLevel = a.loggingOptions.level
	var contentTypeHeader = a._getContentType(request.Header)
	var url string
	if a.loggingOptions.request.includeQueryInPath {
		url = request.RequestURI
	} else {
		url = a._removeQueryParams(request.RequestURI)
	}

	a.logger.Log(
		logLevel,
		fmt.Sprintf("Request %v %v %v", request.Method, url, contentTypeHeader),
		map[string]any{
			"method":      request.Method,
			"url":         url,
			"contentType": contentTypeHeader,
		})

	a._applyLogRequestOptions(logLevel, request)
}

// LogResponse Logs an HTTP response.
func (a *SdkLogger) LogResponse(response *http.Response) {
	var logLevel = a.loggingOptions.level
	var contentTypeHeader = a._getContentType(response.Header)
	var contentLengthHeader = a._getContentLength(response.Header)

	a.logger.Log(
		logLevel,
		fmt.Sprintf("Response %v %v %v", response.StatusCode, contentLengthHeader, contentTypeHeader),
		map[string]any{
			"statusCode":    response.StatusCode,
			"contentLength": contentLengthHeader,
			"contentType":   contentTypeHeader,
		},
	)

	a._applyLogResponseOptions(logLevel, response)
}

func (a *SdkLogger) _applyLogRequestOptions(level Level, request *http.Request) {
	a._applyLogRequestHeaders(
		level,
		request,
		a.loggingOptions.request,
	)

	a._applyLogRequestBody(level, request, a.loggingOptions.request)
}

func (a *SdkLogger) _applyLogRequestHeaders(
	level Level,
	request *http.Request,
	logRequest RequestLoggerConfiguration) {

	logHeaders := logRequest.headers
	headersToInclude := logRequest.includeHeaders
	headersToExclude := logRequest.excludeHeaders
	headersToWhitelist := logRequest.whitelistHeaders

	if logHeaders {
		var headersToLog = a._extractHeadersToLog(
			headersToInclude,
			headersToExclude,
			headersToWhitelist,
			request.Header,
		)

		a.logger.Log(
			level,
			fmt.Sprintf("Request headers %v", headersToLog),
			map[string]any{"headers": headersToLog},
		)
	}
}

func (a *SdkLogger) _applyLogRequestBody(
	level Level,
	request *http.Request,
	logRequest RequestLoggerConfiguration) {

	if logRequest.body {
		a.logger.Log(level, fmt.Sprintf("Request body %v", request.Body),
			map[string]any{"body": request.Body},
		)
	}
}

func (a *SdkLogger) _applyLogResponseOptions(level Level, response *http.Response) {
	a._applyLogResponseHeaders(
		level,
		response,
		a.loggingOptions.response,
	)

	a._applyLogResponseBody(
		level,
		response,
		a.loggingOptions.response,
	)
}

func (a *SdkLogger) _applyLogResponseHeaders(
	level Level,
	response *http.Response,
	logResponse MessageLoggerConfiguration) {

	logHeaders := logResponse.headers
	headersToInclude := logResponse.includeHeaders
	headersToExclude := logResponse.excludeHeaders
	headersToWhitelist := logResponse.whitelistHeaders

	if logHeaders {
		var headersToLog = a._extractHeadersToLog(
			headersToInclude,
			headersToExclude,
			headersToWhitelist,
			response.Header,
		)

		a.logger.Log(level, fmt.Sprintf("Response headers %v", headersToLog),
			map[string]any{"headers": headersToLog},
		)
	}
}

func (a *SdkLogger) _applyLogResponseBody(
	level Level,
	response *http.Response,
	logResponse MessageLoggerConfiguration) {

	if logResponse.body {
		a.logger.Log(level, fmt.Sprintf("Response body %v", response.Body),
			map[string]any{"body": response.Body},
		)
	}
}

const CONTENT_TYPE_HEADER = "content-type"
const CONTENT_LENGTH_HEADER = "content-length"

func (a *SdkLogger) _getContentType(headers http.Header) string {
	var contentType string = ""
	if len(headers) > 0 {
		contentType = headers.Get(CONTENT_TYPE_HEADER)
	}
	return contentType
}

func (a *SdkLogger) _getContentLength(headers http.Header) string {
	var contentLength string = ""
	if len(headers) > 0 {
		contentLength = headers.Get(CONTENT_LENGTH_HEADER)
	}
	return contentLength
}

func (a *SdkLogger) _removeQueryParams(url string) string {
	if strIndex := strings.Index(url, "?"); strIndex != -1 {
		return url[:strIndex]
	}
	return url
}

func (a *SdkLogger) _extractHeadersToLog(
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

func (a *SdkLogger) _includeHeadersToLog(
	headers, filteredHeaders http.Header,
	headersToInclude []string) http.Header {
	// Filter headers based on the keys specified in includeHeaders
	for _, name := range headersToInclude {
		val, ok := headers[name]
		if len(val) > 0 && ok {
			filteredHeaders[name] = val
		}
	}
	return filteredHeaders
}

func (a *SdkLogger) _excludeHeadersToLog(
	headers, filteredHeaders http.Header,
	headersToExclude []string) http.Header {
	// Filter headers based on the keys specified in excludeHeaders
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

func (a *SdkLogger) _maskSenstiveHeaders(
	headers http.Header,
	headersToWhitelist []string) http.Header {

	if a.loggingOptions.maskSensitiveHeaders {
		for key := range headers {
			val := headers.Get(key)
			headers.Set(key, a._maskIfSenstiveHeader(key, val, headersToWhitelist))
		}
	}
	return headers
}

func (a *SdkLogger) _maskIfSenstiveHeader(
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
