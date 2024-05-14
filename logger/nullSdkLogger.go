package logger

import (
	"net/http"
)

// NullSdkLogger represents implementation for SdkLoggerInterface, implementing methods to log HTTP requests and responses.
type NullSdkLogger struct{}

// LogRequest request Logs an HTTP request.
func (a NullSdkLogger) LogRequest(_request *http.Request) {}

// LogResponse Logs an HTTP response.
func (a NullSdkLogger) LogResponse(_response *http.Response) {}
