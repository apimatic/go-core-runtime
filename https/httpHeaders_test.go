package https

import (
	"reflect"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	headers := map[string]string{}
	SetHeaders(headers, "sdk", "go")
	expected := map[string]string{"sdk": "go"}

	if !reflect.DeepEqual(headers, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, headers)
	}
}

func TestMergeHeaders(t *testing.T) {
	headers := map[string]string{"sdk": "go"}
	headersToMerge := map[string]string{"library": "core"}
	MergeHeaders(headers, headersToMerge)
	expected := map[string]string{"sdk": "go", "library": "core"}

	if !reflect.DeepEqual(headers, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, headers)
	}
}

func TestMergeHeadersEmptyMap(t *testing.T) {
	headers := map[string]string{}
	headersToMerge := map[string]string{"library": "core"}
	MergeHeaders(headers, headersToMerge)
	expected := map[string]string{"library": "core"}

	if !reflect.DeepEqual(headers, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, headers)
	}
}
