package logger

import (
	"bytes"
	"io"
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

// NullSdkLogger represents implementation for SdkLoggerInterface, implementing methods to log HTTP requests and responses.
type NullSdkLogger struct{}

// LogRequest request Logs an HTTP request.
func (a NullSdkLogger) LogRequest(*http.Request) {}

// LogResponse Logs an HTTP response.
func (a NullSdkLogger) LogResponse(*http.Response) {}

// SdkLogger represents implementation for SdkLoggerInterface, providing methods to log HTTP requests and responses.
type SdkLogger struct {
	loggingOptions LoggerConfiguration
	logger         LoggerInterface
}

// NewSdkLogger Constructs a new instance of SdkLogger or NullSdkLogger.
func NewSdkLogger(loggingOpt LoggerConfiguration) SdkLoggerInterface {
	if loggingOpt.isValid() {
		return &SdkLogger{
			loggingOptions: loggingOpt,
			logger:         loggingOpt.logger,
		}
	} else {
		return NullSdkLogger{}
	}
}

// LogRequest request Logs an HTTP request.
func (a *SdkLogger) LogRequest(request *http.Request) {
	level := a.loggingOptions.level
	url := request.URL.String()
	if !a.loggingOptions.request.includeQueryInPath {
		url = a._removeQueryParams(url)
	}

	a.logger.Log(
		level,
		"Request %{method} %{url} %{contentType}",
		map[string]any{
			"method":      request.Method,
			"url":         url,
			"contentType": a._getContentType(request.Header),
		})

	a._applyLogRequestOptions(level, request)
}

// LogResponse Logs an HTTP response.
func (a *SdkLogger) LogResponse(response *http.Response) {

	level := a.loggingOptions.level
	a.logger.Log(
		level,
		"Response %{statusCode} %{contentLength} %{contentType}",
		map[string]any{
			"statusCode":    response.StatusCode,
			"contentLength": a._getContentLength(response.Header),
			"contentType":   a._getContentType(response.Header),
		},
	)

	a._applyLogResponseOptions(level, response)
}

func (a *SdkLogger) _applyLogRequestOptions(level Level, request *http.Request) {

	logOp := a.loggingOptions.request
	if logOp.headers {
		a.logger.Log(level, "Request headers %{headers}",
			map[string]any{"headers": a._extractHeadersToLog(
				logOp.includeHeaders,
				logOp.excludeHeaders,
				logOp.whitelistHeaders,
				request.Header,
			)},
		)
	}

	if logOp.body {
		if request.Body == nil {
			a.logger.Log(level, "Request body %{body}", map[string]any{"body": nil})
			return
		}
		bodyBytes, err := io.ReadAll(request.Body)
		request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if err != nil {
			a.logger.Log(level, "Error reading request body %{body}", map[string]any{"body": err})
		} else {
			a.logger.Log(level, "Request body %{body}", map[string]any{"body": string(bodyBytes)})
		}
	}
}

func (a *SdkLogger) _applyLogResponseOptions(level Level, response *http.Response) {

	logOp := a.loggingOptions.response
	if logOp.headers {
		a.logger.Log(level, "Response headers %{headers}",
			map[string]any{"headers": a._extractHeadersToLog(
				logOp.includeHeaders,
				logOp.excludeHeaders,
				logOp.whitelistHeaders,
				response.Header,
			)},
		)
	}

	if logOp.body {
		if response.Body == http.NoBody {
			a.logger.Log(level, "Response body %{body}", map[string]any{"body": nil})
			return
		}
		bodyBytes, err := io.ReadAll(response.Body)
		response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if err != nil {
			a.logger.Log(level, "Error reading response body %{body}", map[string]any{"body": err})
		} else {
			a.logger.Log(level, "Response body %{body}", map[string]any{"body": string(bodyBytes)})
		}
	}
}

const CONTENT_TYPE_HEADER = "content-type"
const CONTENT_LENGTH_HEADER = "content-length"

func (a *SdkLogger) _getContentType(headers http.Header) string {
	var contentType = ""
	if len(headers) > 0 {
		contentType = headers.Get(CONTENT_TYPE_HEADER)
	}
	return contentType
}

func (a *SdkLogger) _getContentLength(headers http.Header) string {
	var contentLength = ""
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
		for key, values := range headers {
			filteredHeaders[key] = values
		}
	}

	return a._maskSensitiveHeaders(filteredHeaders, headersToWhitelist)
}

func (a *SdkLogger) _includeHeadersToLog(
	headers, filteredHeaders http.Header,
	headersToInclude []string) http.Header {
	// Filter headers based on the keys specified in includeHeaders
	for _, name := range headersToInclude {
		nameLower := strings.ToLower(name)
		for headerKey, headerVal := range headers {
			if strings.ToLower(headerKey) == nameLower {
				filteredHeaders[headerKey] = headerVal
			}
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
		if strings.EqualFold(name, key) {
			return true
		}
	}
	return false
}

func (a *SdkLogger) _maskSensitiveHeaders(
	headers http.Header,
	headersToWhitelist []string) http.Header {

	if a.loggingOptions.maskSensitiveHeaders {
		for key := range headers {
			val := headers.Get(key)
			headers.Set(key, a._maskIfSensitiveHeader(key, val, headersToWhitelist))
		}
	}
	return headers
}

func (a *SdkLogger) _maskIfSensitiveHeader(
	name, value string,
	headersToWhiteList []string) string {
	var nonSensitiveHeaders = []string{
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
