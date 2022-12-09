package testHelper

import "testing"

func TestCheckResponseStatusCode(t *testing.T) {
	result := 202
	expected := 202
	CheckResponseStatusCode(t, result, expected)
}

func TestCheckResponseStatusCodeError(t *testing.T) {
	result := 202
	expected := 200
	CheckResponseStatusCode(&testing.T{}, result, expected)
}

func TestCheckResponseStatusRange(t *testing.T) {
	result := 208
	expectedLowerLimit := 200
	expectedUpperLimit := 208
	CheckResponseStatusRange(t, result, expectedLowerLimit, expectedUpperLimit)
}

func TestCheckResponseStatusRangeError(t *testing.T) {
	result := 212
	expectedLowerLimit := 200
	expectedUpperLimit := 208
	CheckResponseStatusRange(&testing.T{}, result, expectedLowerLimit, expectedUpperLimit)
}
