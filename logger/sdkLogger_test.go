package logger_test

import (
	"context"
	"fmt"
	"github.com/apimatic/go-core-runtime/https"
	. "github.com/apimatic/go-core-runtime/logger"
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
	request.AppendPath("/integer")
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

type fmtLogger struct{}

// Logs a message to the console with the specified log level.
func (c fmtLogger) Log(level Level, message string, params map[string]any) {
	fmt.Printf("Level : %v,\t Message : %v,\t Params : %v", level, message, params)
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
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLogger(fmtLogger{}),
	)))
	_callRequestAsJson(t)
}

func TestSDKLoggerWithCustomLoggerCustomConfig(t *testing.T) {
	request.Logger(NewSdkLogger(NewLoggerConfiguration(
		WithLevel("debug"),
		WithLogger(fmtLogger{}),
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
}
