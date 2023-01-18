package https

import (
	"net/http"
	"reflect"
	"testing"
)

func TestApiResponse(t *testing.T) {
	apiResponse := ApiResponse[string]{
		Data:     "This is data",
		Response: &http.Response{}}
	expected := NewApiResponse("This is data", &http.Response{})
	if !reflect.DeepEqual(apiResponse, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, apiResponse)
	}
}
