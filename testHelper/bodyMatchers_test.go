package testHelper

import (
	"testing"
)

func TestNativeBodyMatcherNumber(t *testing.T) {
	expected := `4`
	var result int = 4
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherPrecision(t *testing.T) {
	expected := `4.11`
	var result float32 = 4.11
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherLong(t *testing.T) {
	expected := `411111111111111111`
	var result int64 = 411111111111111111
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherBoolean(t *testing.T) {
	expected := `true`
	var result bool = true
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherStringSlice(t *testing.T) {
	expected := `["Tuesday", "Saturday", "Wednesday", "Monday", "Sunday"]`
	var result []string = []string{
		"Tuesday", "Saturday", "Wednesday", "Monday", "Sunday",
	}
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherIntSlice(t *testing.T) {
	expected := `[1,2,3,4,5]`
	var result []int = []int{
		1, 2, 3, 4, 5,
	}
	NativeBodyMatcher(t, expected, result)
}

type Object struct {
	Key1 string
	Key2 int
}

func TestKeysAndValuesBodyMatcher(t *testing.T) {
	expected := `{"Key1":"value","Key2":1}`
	result := Object{
		Key1: "value",
		Key2: 1,
	}
	KeysAndValuesBodyMatcher(t, expected, result, false, false)
}

func TestKeysAndValuesBodyMatcherEmpty(t *testing.T) {
	expected := `{}`
	KeysAndValuesBodyMatcher(t, expected, nil, false, false)
}
