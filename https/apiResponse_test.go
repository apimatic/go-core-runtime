package https

import (
	"net/http"
	"reflect"
	"testing"
)

func TestApiResponse(t *testing.T) {
	apiResponse := ApiResponse[string] {
		Data: "This is error body",
		Response: &http.Response{}}
	expected := NewApiResponse("This is error body", &http.Response{})
	if !reflect.DeepEqual(apiResponse, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, apiResponse)
	}
}
