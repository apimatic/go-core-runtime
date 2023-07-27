package testHelper

import (
	"net/http"
	"testing"
)

// TestHeader represents a test header with Name, Value, and CheckValue fields.
type TestHeader struct {
	CheckValue bool   `json:"CheckValue"`
	Name       string `json:"Name"`
	Value      string `json:"Value"`
}

// NewTestHeader creates and returns a new TestHeader with the given values.
func NewTestHeader(checkValue bool, name, value string) TestHeader {
	return TestHeader{
		CheckValue: checkValue,
		Name:       name,
		Value:      value,
	}
}

// CheckResponseHeaders checks if the responseHeaders contain the expected headers specified in expectedHeadersList.
// If allowExtraHeaders is set to false, it also verifies that no extra headers are present in the responseHeaders.
func CheckResponseHeaders(t *testing.T, responseHeaders http.Header, expectedHeadersList []TestHeader, allowExtraHeaders bool) {
	for _, expectedHeader := range expectedHeadersList {
		respValue := responseHeaders.Get(expectedHeader.Name)
		if respValue == "" {
			t.Errorf("expected header '%v' does not exists in response", expectedHeader.Name)
			break
		} else if expectedHeader.CheckValue && respValue != expectedHeader.Value {
			t.Errorf("response does not contains same value of expected header '%v'", expectedHeader.Name)
			break
		}
	}

	if !allowExtraHeaders && len(responseHeaders) != len(expectedHeadersList) {
		t.Errorf("response contains other headers than those listed in the expected headers list")
	}
}
