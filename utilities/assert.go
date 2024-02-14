package utilities

import "testing"

// Asserts that both values are equal.
// This may not apply to all types.
func AssertEquals[T comparable](t *testing.T, expected, actual T) {
	if expected == actual {
		return
	}
	t.Errorf("\n----------EXPECTED----------\n%v\n------------GOT------------\n%v", expected, actual)
}

// Asserts that an error is returned (i.e. not nil).
func AssertError(t *testing.T, err error) {
	if err != nil {
		return
	}
	t.Errorf("Expected an error")
}

// Asserts that no error is returned (i.e. nil).
func AssertNoError(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatalf("Unexpected Error: %s", err.Error())
}
