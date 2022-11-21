package https

import (
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
	client := NewHttpClient()
	callBuilder := CreateCallBuilderFactory(
		func(server string) string { return "https://apimatic-go.free.beeceptor.com" }, nil, client)

	req := callBuilder("POST", "/interceptors")
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
	respBody, response, err := req.CallAsText()
	if err != nil {
		t.Errorf("Error in CallAsText: %v", err)
	}

	if response.StatusCode != 200 || respBody != "Success" {
		t.Errorf("Interceptors not working!")
	}
}
