package apiError

import (
	"reflect"
	"testing"
)

func TestNewApiError(t *testing.T) {
	expected := &ApiError{
		StatusCode: 500,
		Body:       "Server Error",
	}
	result := NewApiError(500, "Server Error")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestErrorMethod(t *testing.T) {
	expected := "ApiError occured Server Error"
	result := NewApiError(500, "Server Error")
	if !reflect.DeepEqual(result.Error(), expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Error())
	}
}
