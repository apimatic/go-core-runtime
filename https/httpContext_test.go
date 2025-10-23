package https

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestAddQuery(t *testing.T) {

	callBuilder := GetCallBuilder(ctx, "GET", "", nil)
	request, err := callBuilder.toRequest()
	AddQuery(request, "param", "query")

	if request.URL.RawQuery != "param=query" || err != nil {
		t.Errorf("Failed:\nExpected query param missing")
	}
}

func TestGetResponseBody(t *testing.T) {
	bodyBytes := []byte(`{"invalidJson"}`)
	respBody := io.NopCloser(bytes.NewReader(bodyBytes))
	ctx := HttpContext{
		Response: &http.Response{
			Body: respBody,
		},
	}

	newBodyBytes, err := ctx.GetResponseBody()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(newBodyBytes, bodyBytes) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", bodyBytes, newBodyBytes)
	}
}

func TestReadRequestBody(t *testing.T) {
	bodyBytes := []byte(`{"invalidJson"}`)
	respBody := io.NopCloser(bytes.NewReader(bodyBytes))
	ctx := HttpContext{
		Request: &http.Request{
			Body: respBody,
		},
	}

	newBodyBytes, err := ReadRequestBody(ctx.Request)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(newBodyBytes, bodyBytes) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", bodyBytes, newBodyBytes)
	}
}
