package https

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"strconv"
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

func TestHttpClientInterfaceExecuteFailure(t *testing.T) {
	req, err := FromFastHttpRequest(
		[]byte("INVALID_METHOD"),
		"://invalid-url",
		nil,
		nil,
		[]byte{},
		nil,
	)

	if err == nil || req != nil {
		t.Errorf("Failed: Response not okay!\n %v", req)
	}
}

func TestHttpClientInterfaceExecute(t *testing.T) {
	fastReq, _ := http.NewRequest(
		http.MethodGet,
		"https://apimatic-go.free.beeceptor.com",
		io.NopCloser(bytes.NewBuffer([]byte("This is data"))))
	bodyBytes, _ := ReadRequestBody(fastReq)

	fastReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fastReq.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))

	fastReqVisitAll := func(setFunc func(k, v []byte)) {
		for k, v := range fastReq.Header {
			for _, hv := range v {
				setFunc([]byte(k), []byte(hv))
			}
		}
	}

	fastReq.AddCookie(&http.Cookie{Name: API_KEY, Value: "12345"})
	fastReq.AddCookie(&http.Cookie{Name: API_TOKEN, Value: "token-123456"})

	fastReqVisitAllCookie := func(setFunc func(key, value []byte)) {
		for _, cookie := range fastReq.Cookies() {
			setFunc([]byte(cookie.Name), []byte(cookie.Value))
		}
	}

	req, err := FromFastHttpRequest(
		[]byte(fastReq.Method),
		fastReq.URL.String(),
		bodyBytes,
		fastReqVisitAll,
		[]byte(fastReq.RequestURI),
		fastReqVisitAllCookie,
	)

	if err != nil || req == nil {
		t.Errorf("Failed: Response not okay!\n %v", req)
	}
}
