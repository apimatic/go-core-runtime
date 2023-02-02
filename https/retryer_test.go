package https

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func getRetryConfiguration() RetryConfiguration {
	return NewRetryConfiguration(
		WithBackoffFactor(5),
		WithMaxRetryAttempts(11),
		WithRetryOnTimeout(false),
		WithHttpMethodsToRetry([]string{"DELETE"}),
		WithHttpStatusCodesToRetry([]int64{429}),
		WithRetryInterval(5),
		WithMaximumRetryWaitTime(50),
	)
}

func TestRetryConfigurationOptions(t *testing.T) {
	got := getRetryConfiguration()

	if got.BackoffFactor() != 5 ||
		got.MaxRetryAttempts() != 11 ||
		got.RetryOnTimeout() != false ||
		!reflect.DeepEqual(got.HttpMethodsToRetry(), []string{"DELETE"}) ||
		!reflect.DeepEqual(got.HttpStatusCodesToRetry(), []int64{429}) ||
		got.RetryInterval() != 5 ||
		got.MaximumRetryWaitTime() != 50 {
		t.Errorf("Failed:\nGot: %v", got)
	}
}

func TestRequestRetryOptionDefault(t *testing.T) {
	var retryOptionDefault RequestRetryOption = 0

	if retryOptionDefault.String() != "default" {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", "default", retryOptionDefault.String())
	}
}

func TestRequestRetryOptionEnable(t *testing.T) {
	var retryOptionEnable RequestRetryOption = 1

	if retryOptionEnable.String() != "enable" {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", "enable", retryOptionEnable.String())
	}
}
func TestRequestRetryOptionDisable(t *testing.T) {
	var retryOptionDisable RequestRetryOption = 2

	if retryOptionDisable.String() != "disable" {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", "disable", retryOptionDisable.String())
	}
}

func TestContainsHttpStatusCodesToRetry429(t *testing.T) {
	retryConfig := getRetryConfiguration()
	if !retryConfig.containsHttpStatusCodesToRetry(429) {
		t.Error("Failed:\nRetry Configuration HttpStatusCode has 429 but method returns false.")
	}
}

func TestContainsHttpStatusCodesToRetry200(t *testing.T) {
	retryConfig := getRetryConfiguration()
	if retryConfig.containsHttpStatusCodesToRetry(200) {
		t.Error("Failed:\nRetry Configuration HttpStatusCode does not have 200 but the method returns true.")
	}
}

func TestContainsHttpMethodsToRetryDelete(t *testing.T) {
	retryConfig := getRetryConfiguration()
	if !retryConfig.containsHttpMethodsToRetry(http.MethodDelete) {
		t.Error("Failed:\nRetry Configuration HttpMethodsToRetry has MethodDelete but method returns false.")
	}
}

func TestContainsHttpMethodsToRetryPost(t *testing.T) {
	retryConfig := getRetryConfiguration()
	if retryConfig.containsHttpMethodsToRetry(http.MethodPost) {
		t.Error("Failed:\nRetry Configuration HttpMethodsToRetry does not have MethodPost but the method returns true.")
	}
}

func TestGetRetryAfterInSecondsEmptyHeaders(t *testing.T) {
	headers := http.Header{}
	if getRetryAfterInSeconds(headers) != 0 {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", 0, getRetryAfterInSeconds(headers))
	}
}

func TestGetRetryAfterInSeconds(t *testing.T) {
	headers := http.Header{}
	headers.Add("retry-after", "20")
	if getRetryAfterInSeconds(headers) != 20 {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", 20, getRetryAfterInSeconds(headers))
	}
}

func TestDefaultBackoff(t *testing.T) {
	backoff := defaultBackoff(5, 100, 20, 2, 1)
	if time.Duration(backoff) != time.Duration(20) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", time.Duration(20)*time.Second, time.Duration(backoff)*time.Second)
	}
}

func TestDefaultBackoffZero(t *testing.T) {
	backoff := defaultBackoff(5, 1, 20, 2, 1)
	if time.Duration(backoff) != time.Duration(0) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", time.Duration(0)*time.Second, time.Duration(backoff)*time.Second)
	}
}

func TestShouldRetryEnable(t *testing.T) {
	retryConfig := getRetryConfiguration()
	shouldRetry := retryConfig.ShouldRetry(RequestRetryOption(Enable), "")
	if !shouldRetry {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", shouldRetry, false)
	}
}

func TestShouldRetryDisable(t *testing.T) {
	retryConfig := getRetryConfiguration()
	shouldRetry := retryConfig.ShouldRetry(RequestRetryOption(Disable), "")
	if shouldRetry {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", shouldRetry, true)
	}
}

func TestShouldRetryDefaultFalse(t *testing.T) {
	retryConfig := getRetryConfiguration()
	shouldRetry := retryConfig.ShouldRetry(RequestRetryOption(Default), http.MethodGet)
	if shouldRetry {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", shouldRetry, true)
	}
}

func TestShouldRetryDefaultTrue(t *testing.T) {
	retryConfig := NewRetryConfiguration(
		WithMaxRetryAttempts(11),
		WithHttpMethodsToRetry([]string{"DELETE"}),
		WithHttpStatusCodesToRetry([]int64{429}),
		WithMaximumRetryWaitTime(50),
	)
	shouldRetry := retryConfig.ShouldRetry(RequestRetryOption(Default), http.MethodDelete)
	if !shouldRetry {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", shouldRetry, false)
	}
}
func TestGetRetryWaitTimeTimeoutError(t *testing.T) {
	httpClient := http.Client{Timeout: 1 * time.Nanosecond}
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	_, err := httpClient.Do(req)
	retryConfig := getRetryConfiguration()
	waitTime := retryConfig.GetRetryWaitTime(0, 1, nil, err)
	if waitTime != 0 {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", time.Duration(0)*time.Second, waitTime)
	}
}

func TestGetRetryWaitTimeRetryAfterHeader(t *testing.T) {
	resp := http.Response{
		Header: http.Header{},
	}
	resp.Header.Add("retry-after", "20")
	retryConfig := getRetryConfiguration()
	waitTime := retryConfig.GetRetryWaitTime(0, 1, &resp, nil)
	if waitTime != 0 {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", time.Duration(0)*time.Second, waitTime)
	}
}

func TestGetRetryWaitTimeRetryStatusCode(t *testing.T) {
	resp := http.Response{
		StatusCode: 400,
	}
	retryConfig := getRetryConfiguration()
	waitTime := retryConfig.GetRetryWaitTime(0, 1, &resp, nil)
	if waitTime != 0 {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", time.Duration(0)*time.Second, waitTime)
	}
}
