package https

import (
	"testing"
)

func TestErrorMethod(t *testing.T) {
	expected := "ApiError occured Server Error"
	result := ApiError{
		StatusCode: 500,
		Body:       "Server Error",
	}
	if result.Error() != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Error())
	}
}
