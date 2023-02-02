package https

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewHttpClient(t *testing.T) {
	result := NewHttpClient(NewHttpConfiguration())
	expected := HttpClient{
		httpClientInstance: *http.DefaultClient,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestHttpClientExecute(t *testing.T) {
	client := NewHttpClient(NewHttpConfiguration())
	response, _ := client.Execute(&http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: "apimatic-go.free.beeceptor.com"}})

	if response.StatusCode != 200 {
		t.Errorf("Failed: Response not okay!\n %v", response)
	}
}

func TestHttpClientExecuteError(t *testing.T) {
	client := NewHttpClient(NewHttpConfiguration())
	response, err := client.Execute(&http.Request{})

	if err == nil && response.StatusCode != 200 {
		t.Errorf("Failed: Response not okay!\n %v", response)
	}
}
