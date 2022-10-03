package https

import (
	"net/http"
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	expected := HttpConfiguration{
		Timeout:   0,
		Transport: http.DefaultTransport,
	}
	got := DefaultHttpConfiguration()

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, got)
	}
}
