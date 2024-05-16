package logger_test

import (
	"context"
	"fmt"
	"github.com/apimatic/go-core-runtime/https"
	. "github.com/apimatic/go-core-runtime/logger"
	"reflect"
	"testing"
)

var callBuilder https.CallBuilderFactory
var ctx = context.Background()
var serverUrl = https.GetTestingServer().URL

func init() {

	client := https.NewHttpClient(https.NewHttpConfiguration())
	callBuilder = https.CreateCallBuilderFactory(
		func(server string) string {
			return serverUrl
		},
		nil,
		client,
		https.NewRetryConfiguration(),
		https.Indexed,
	)
}

func _callRequestAsJson(t *testing.T, request https.CallBuilder) {

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
	request := callBuilder(ctx, "GET", "//response/")
	request.AppendPath("/binary")
	request.Logger(NullSdkLogger{})
	_callRequestAsJson(t, request)
}

func TestSDKLoggerWithDefaultConfig(t *testing.T) {
	request := callBuilder(ctx, "GET", "//response/")
	request.AppendPath("/binary")
	request.Logger(NewSdkLogger(NewLoggerConfiguration()))
	_callRequestAsJson(t, request)
}

func TestSDKLoggerWithInvalidConfig(t *testing.T) {
	request := callBuilder(ctx, "GET", "//response/")
	request.AppendPath("/binary")
	request.Logger(NewSdkLogger(LoggerConfiguration{}))
	_callRequestAsJson(t, request)
}

func TestSDKLoggerWithCustomConfig(t *testing.T) {
	request := callBuilder(ctx, "GET", "//response/")
	request.Header("Authorization", "ahsfhafu3264basfcasbk__aksdk")
	request.AppendPath("/binary")
	request.Header("Content-Type", "application/file")
	request.QueryParam("env", "testing")

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
	_callRequestAsJson(t, request)
}

func TestSDKLoggerWithCustomLoggerDefaultConfig(t *testing.T) {
	request := callBuilder(ctx, "GET", "//response/")
	request.AppendPath("/binary")
	logger := &fmtLogger{}
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLogger(logger),
		WithRequestConfiguration(
			WithRequestHeaders(true),
		),
	)))
	_callRequestAsJson(t, request)

	expected := []string{
		"info, Request %{method} %{url} %{contentType}, map[contentType: method:GET url:" + serverUrl + "/response/binary]",
		"info, Request headers %{headers}, map[headers:map[]]",
		"info, Response %{statusCode} %{contentLength} %{contentType}, map[contentLength:45 contentType:text/plain; charset=utf-8 statusCode:200]",
	}
	logger.AssertLogEntries(t, expected...)
}

func TestSDKLoggerWithCustomLoggerDefaultConfigWithHeaders(t *testing.T) {
	request := callBuilder(ctx, "GET", "//response/")
	request.AppendPath("/binary")
	request.Header("Content-Type", "application/file")
	logger := &fmtLogger{}
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLogger(logger),
		WithRequestConfiguration(
			WithRequestHeaders(true),
		),
	)))
	_callRequestAsJson(t, request)

	expected := []string{
		"info, Request %{method} %{url} %{contentType}, map[contentType:application/file method:GET url:" + serverUrl + "/response/binary]",
		"info, Request headers %{headers}, map[headers:map[Content-Type:[application/file]]]",
		"info, Response %{statusCode} %{contentLength} %{contentType}, map[contentLength:45 contentType:text/plain; charset=utf-8 statusCode:200]",
	}
	logger.AssertLogEntries(t, expected...)
}

func TestSDKLoggerWithCustomLoggerCustomConfig(t *testing.T) {
	request := callBuilder(ctx, "GET", "//response/")
	request.Header("Authorization", "ahsfhafu3264basfcasbk__aksdk")
	request.AppendPath("/binary")
	request.Header("Content-Type", "application/file")
	request.QueryParam("env", "testing")

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
	_callRequestAsJson(t, request)

	expected := []string{
		"debug, Request %{method} %{url} %{contentType}, map[contentType:application/file method:GET url:" + serverUrl + "/response/binary]",
		"debug, Request headers %{headers}, map[headers:map[Authorization:[**Redacted**] Content-Type:[application/file]]]",
		"debug, Request body %{body}, map[body:null]",
		"debug, Response %{statusCode} %{contentLength} %{contentType}, map[contentLength:45 contentType:text/plain; charset=utf-8 statusCode:200]",
		"debug, Response headers %{headers}, map[headers:map[Content-Type:[text/plain; charset=utf-8]]]",
		"debug, Response body %{body}, map[body:{}]",
	}
	logger.AssertLogEntries(t, expected...)
}
