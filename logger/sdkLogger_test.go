package logger_test

import (
	"context"
	"fmt"
	"github.com/apimatic/go-core-runtime/https"
	. "github.com/apimatic/go-core-runtime/logger"
	"reflect"
	"testing"
)

var request https.CallBuilder

func init() {
	ctx := context.Background()
	client := https.NewHttpClient(https.NewHttpConfiguration())
	callBuilder := https.CreateCallBuilderFactory(
		func(server string) string {
			return https.GetTestingServer().URL
		},
		nil,
		client,
		https.NewRetryConfiguration(),
		https.Indexed,
	)

	request = callBuilder(ctx, "GET", "//response/")
	request.AppendPath("/binary")
	request.Header("Content-Type", "application/file")

}

func _callRequestAsJson(t *testing.T) {
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}
	expected := 200
	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

type fmtLogger struct {
	entries []string
}

func (c *fmtLogger) AssertLogEntries(t *testing.T, expected ...string) {
	if !reflect.DeepEqual(c.entries, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, c.entries)
	}
}

// Logs a message to the console with the specified log level.
func (c *fmtLogger) Log(level Level, message string, params map[string]any) {
	c.entries = append(c.entries, fmt.Sprintf("%v, %v, %v", level, message, params))
}

func TestNullSDKLogger(t *testing.T) {
	request.Logger(NullSdkLogger{})
	_callRequestAsJson(t)
}

func TestSDKLoggerWithDefaultConfig(t *testing.T) {
	request.Logger(NewSdkLogger(NewLoggerConfiguration()))
	_callRequestAsJson(t)
}

func TestSDKLoggerWithCustomConfig(t *testing.T) {
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLevel("debug"),
		WithMaskSensitiveHeaders(true),
		WithRequestConfiguration(
			WithRequestBody(true),
			WithRequestHeaders(true),
			WithIncludeQueryInPath(true),
			WithIncludeRequestHeaders("Content-Type", "Content-Encoding"),
			WithWhitelistRequestHeaders("Authorization"),
		),
		WithResponseConfiguration(
			WithResponseBody(true),
			WithExcludeResponseHeaders("X-Powered-By"),
		),
	)))
	_callRequestAsJson(t)
}

func TestSDKLoggerWithCustomLoggerDefaultConfig(t *testing.T) {
	logger := &fmtLogger{}
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLogger(logger),
	)))
	_callRequestAsJson(t)

	expected := []string{
		"info, Request %{method} %{url} %{contentType}, map[contentType: method:GET url:]",
		"info, Response %{statusCode} %{contentLength} %{contentType}, map[contentLength:45 contentType:text/plain; charset=utf-8 statusCode:200]",
	}
	logger.AssertLogEntries(t, expected...)
}

func TestSDKLoggerWithCustomLoggerCustomConfig(t *testing.T) {
	logger := &fmtLogger{}
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLevel("debug"),
		WithLogger(logger),
		WithMaskSensitiveHeaders(true),
		WithRequestConfiguration(
			WithRequestBody(true),
			WithRequestHeaders(true),
			WithExcludeRequestHeaders("X-Powered-By",
				"range",
				"referer",
				"retry-after"),
		),
		WithResponseConfiguration(
			WithResponseBody(true),
			WithResponseHeaders(true),
			WithIncludeResponseHeaders("Content-Type", "Content-Encoding"),
			WithWhitelistResponseHeaders("Authorization"),
		),
	)))
	_callRequestAsJson(t)

	expected := []string{
		"debug, Request %{method} %{url} %{contentType}, map[contentType: method:GET url:]",
		"debug, Request headers %{headers}, map[headers:map[]]",
		"debug, Request body %{body}, map[body:null]",
		"debug, Response %{statusCode} %{contentLength} %{contentType}, map[contentLength:45 contentType:text/plain; charset=utf-8 statusCode:200]",
		"debug, Response headers %{headers}, map[headers:map[Content-Type:[text/plain; charset=utf-8]]]",
		"debug, Response body %{body}, map[body:{}]",
	}
	logger.AssertLogEntries(t, expected...)
}
