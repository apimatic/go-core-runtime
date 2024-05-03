package https

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestPassThroughInterceptor(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "www.google.com", nil)

	next := func(request *http.Request) HttpContext {
		request.Header.Set("apimatic", "go-sdk")
		return HttpContext{Request: request}
	}
	result := PassThroughInterceptor(request, next)
	expected, _ := http.NewRequest(http.MethodGet, "www.google.com", nil)
	expected.Header.Add("apimatic", "go-sdk")

	if !reflect.DeepEqual(result.Request, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Request)
	}
}

func TestPassThroughInterceptorNilNext(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "www.google.com", nil)

	next := func(request *http.Request) HttpContext {
		return HttpContext{Request: request}
	}
	result := PassThroughInterceptor(request, next)
	expected, _ := http.NewRequest(http.MethodGet, "www.google.com", nil)

	if !reflect.DeepEqual(result.Request, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Request)
	}
}

func TestCallHttpInterceptors(t *testing.T) {
	client := NewHttpClient(NewHttpConfiguration())
	callBuilder := CreateCallBuilderFactory(
		func(server string) string { return "https://apimatic-goo.free.beeceptor.com" }, nil, client, NewRetryConfiguration(), Indexed, &ApiLogger{})

	req := callBuilder(context.Background(), "POST", "/interceptors")
	interceptor1 := func(request *http.Request, next HttpCallExecutor) HttpContext {
		request.Header.Set("apimatic", "go-sdk")
		return next(request)
	}
	interceptor2 := func(request *http.Request, next HttpCallExecutor) HttpContext {
		request.Header.Set("library", "core")
		return next(request)
	}
	req.intercept(interceptor1)
	req.intercept(interceptor2)
	_, response, err := req.CallAsText()
	if err != nil {
		t.Errorf("Error in CallAsText: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Interceptors not working!")
	}
}
