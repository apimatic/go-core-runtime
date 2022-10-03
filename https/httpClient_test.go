package https

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	result := NewHttpClient()
	expected := HttpClient{
		httpClientInstance: *http.DefaultClient,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestNewHttpClientWithSingleOption(t *testing.T) {
	result := NewHttpClient(WithTimeout(1))
	expected := HttpClient{
		httpClientInstance: *http.DefaultClient,
	}
	expected.httpClientInstance.Timeout = time.Duration(1 * float64(time.Second))

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestNewHttpClientWithOptions(t *testing.T) {
	result := NewHttpClient(
		WithTimeout(1),
		WithTransport(http.DefaultTransport),
	)
	expected := HttpClient{
		httpClientInstance: *http.DefaultClient,
	}
	expected.httpClientInstance.Timeout = time.Duration(1 * float64(time.Second))
	expected.httpClientInstance.Transport = http.DefaultTransport

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestHttpClientExecute(t *testing.T) {
	client := NewHttpClient()
	response := client.Execute(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: "apimatic-go.free.beeceptor.com"}})

	if response.StatusCode != 200 {
		t.Errorf("Failed: Response not okay!\n %v", response)
	}
}

func TestHttpClientExecuteError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should panic because request is empty.")
		}
	}()

	client := NewHttpClient()
	response := client.Execute(&http.Request{})

	if response.StatusCode != 200 {
		t.Errorf("Failed: Response not okay!\n %v", response)
	}
}
