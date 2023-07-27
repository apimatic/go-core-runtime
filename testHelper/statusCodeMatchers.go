package testHelper

import (
	"testing"
)

// CheckResponseStatusCode checks if the actual response status code matches the expected status code.
// If they do not match, it reports an error with the actual and expected status codes.
func CheckResponseStatusCode(t *testing.T, respStatus int, expectedStatus int) {
	if respStatus != expectedStatus {
		t.Errorf("got %v but expected %v", respStatus, expectedStatus)
	}
}

// CheckResponseStatusRange checks if the actual response status code falls within the expected range (inclusive).
// If the actual status code is not within the range, it reports an error with the actual status code and the expected range.
func CheckResponseStatusRange(t *testing.T, respStatus int, expectedLowerLimit int, expectedUpperLimit int) {
	if respStatus < expectedLowerLimit || respStatus > expectedUpperLimit {
		t.Errorf("got %v but expected between %v and %v", respStatus, expectedLowerLimit, expectedUpperLimit)
	}
}
