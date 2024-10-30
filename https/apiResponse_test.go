package https_test

import (
	"github.com/apimatic/go-core-runtime/https"
	"net/http"
	"reflect"
	"testing"
)

func TestApiResponse(t *testing.T) {
	apiResponse := https.ApiResponse[string]{
		Data:     "This is data",
		Response: &http.Response{}}
	expected := https.NewApiResponse("This is data", &http.Response{})
	if !reflect.DeepEqual(apiResponse, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, apiResponse)
	}
}
