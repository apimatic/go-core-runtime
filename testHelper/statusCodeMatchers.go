package testHelper

import (
	"testing"
)

func CheckResponseStatusCode(t *testing.T, respStatus int, expectedStatus int) {
	if respStatus != expectedStatus {
		t.Errorf("got %v but expected %v", respStatus, expectedStatus)
	}
}

func CheckResponseStatusRange(t *testing.T, respStatus int, expectedLowerLimit int, expectedUpperLimit int) {
	if respStatus < expectedLowerLimit || respStatus > expectedUpperLimit {
		t.Errorf("got %v but expected between %v and %v", respStatus, expectedLowerLimit, expectedUpperLimit)
	}
}
