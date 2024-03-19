// Package testHelper provides helper functions for testing purposes.
// Copyright (c) APIMatic. All rights reserved.
package testHelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

// RawBodyMatcher checks if the expectedBody is contained within the JSON response body.
func RawBodyMatcher[T any](test *testing.T, expectedBody string, responseObject T) {
	responseBytes, _ := json.Marshal(&responseObject)
	responseBody := string(responseBytes)

	if !strings.Contains(responseBody, expectedBody) {
		test.Errorf("got \n%v \nbut expected %v", responseBody, expectedBody)
	}
}

// NativeBodyMatcher compares the JSON response body with the expected JSON body.
func NativeBodyMatcher[T any](test *testing.T, expectedBody string, responseObject T) {
	responseBytes, _ := json.Marshal(&responseObject)
	var expected, response any
	expectedError := json.Unmarshal([]byte(expectedBody), &expected)
	responseError := json.Unmarshal(responseBytes, &response)

	if expectedError != nil || responseError != nil {
		test.Error("error while unmarshalling for comparison")
	}

	if !reflect.DeepEqual(response, expected) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

// KeysBodyMatcher compares the JSON response body with the expected JSON body using keys only.
// The responseObject and expectedBody should have the same keys.
func KeysBodyMatcher[T any](test *testing.T, expectedBody string, responseObject T, checkArrayCount, checkArrayOrder bool) {
	responseBytes, _ := json.Marshal(&responseObject)

	if !matchKeysAndValues(responseBytes, []byte(expectedBody), checkArrayCount, checkArrayOrder, false) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

// KeysAndValuesBodyMatcher compares the JSON response body with the expected JSON body using keys and values.
// The responseObject and expectedBody should have the same keys and their corresponding values should be equal.
func KeysAndValuesBodyMatcher[T any](test *testing.T, expectedBody string, responseObject T, checkArrayCount, checkArrayOrder bool) {
	responseBytes, _ := json.Marshal(&responseObject)

	if !matchKeysAndValues(responseBytes, []byte(expectedBody), checkArrayCount, checkArrayOrder, true) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

// matchKeysAndValues is a helper function used by KeysBodyMatcher and KeysAndValuesBodyMatcher
// to compare the bytes for keys and values.
func matchKeysAndValues(actualBytes, expectedBytes []byte, checkArrayCount, checkArrayOrder, checkValues bool) bool {
	actual, expected, ok := extractAsMaps(actualBytes, expectedBytes)
	if !ok {
		return false
	}
	return matchKeysAndValuesAsMap(actual, expected, checkArrayCount, checkArrayOrder, checkValues)
}

// matchKeysAndValuesAsMap is a helper function used by matchKeysAndValues
// to compare the JSON keys and values recursively.
func matchKeysAndValuesAsMap(actual, expected map[string]any, checkArrayCount, checkArrayOrder, checkValues bool) bool {
	if checkArrayCount && len(expected) != len(actual) {
		return false
	}
	for key, value := range expected {
		responseValue := actual[key]
		if reflect.ValueOf(responseValue).Kind() == reflect.Map {
			if reflect.ValueOf(value).Kind() != reflect.Map {
				return false
			}
			responseSubMap := responseValue.(map[string]any)
			expectedSubMap := value.(map[string]any)
			if !matchKeysAndValuesAsMap(responseSubMap, expectedSubMap, checkArrayCount, checkArrayOrder, checkValues) {
				return false
			}
		} else if checkValues && !reflect.DeepEqual(responseValue, value) {
			return false
		}
	}
	return true
}

// extractAsMaps converts the bytes into maps, if the bytes yield as array,
// then the array indexes will be used as map keys.
func extractAsMaps(actualBytes, expectedBytes []byte) (map[string]any, map[string]any, bool) {
	var actual, expected = make(map[string]any), make(map[string]any)
	var actualErr, expectedErr error

	if len(expectedBytes) > 0 && expectedBytes[0] == '[' {
		var actualArray, expectedArray []any
		actualErr = json.Unmarshal(actualBytes, &actualArray)
		expectedErr = json.Unmarshal(expectedBytes, &expectedArray)
		if actualErr != nil || expectedErr != nil {
			return nil, nil, false
		}
		for i, v := range actualArray {
			actual[fmt.Sprint(i)] = v
		}
		for i, v := range expectedArray {
			expected[fmt.Sprint(i)] = v
		}
	} else {
		actualErr = json.Unmarshal(actualBytes, &actual)
		expectedErr = json.Unmarshal(expectedBytes, &expected)
	}

	if actualErr != nil || expectedErr != nil {
		return nil, nil, false
	}
	return actual, expected, true
}

// IsSameAsFile checks if the responseFileBytes is the same as the content of the file fetched from the expectedFileURL.
func IsSameAsFile(test *testing.T, expectedFileURL string, responseFileBytes []byte) {
	expectedFile, err := https.GetFile(expectedFileURL)
	if err != nil {
		test.Errorf("Cannot get the file: %v", err)
	}
	IsSameInputBytes(test, expectedFile.File, responseFileBytes)
}

// IsSameInputBytes checks if the receivedBytes are equal to the expectedBytes.
func IsSameInputBytes(test *testing.T, expectedBytes []byte, receivedBytes []byte) {
	if !reflect.DeepEqual(expectedBytes, receivedBytes) {
		test.Error("Received bytes do not match the bytes expected")
	}
}

// SliceToCommaSeparatedString converts a slice to a comma-separated string representation.
func SliceToCommaSeparatedString(slice any) string {
	return strings.Join(strings.Split(fmt.Sprint(slice), " "), ",")
}
