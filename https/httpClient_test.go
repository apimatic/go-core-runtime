package https

import (
	"net/http"
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
	request, _ := http.NewRequest(http.MethodGet, "https://apimatic-go.free.beeceptor.com", nil)
	response, _ := client.Execute(request)

	if response == nil || response.StatusCode != 200 {
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
