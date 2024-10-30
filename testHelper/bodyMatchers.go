// Package testHelper provides helper functions for testing purposes.
// Copyright (c) APIMatic. All rights reserved.
package testHelper

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

func setTestingError(test *testing.T, responseArg any, expectedArg any) {
	test.Errorf("got \n%v \nbut expected %v", responseArg, expectedArg)
}

func setTestResponseError(test *testing.T, responseErr error) {
	test.Errorf("Invalid response data: %v", responseErr)
}

// RawBodyMatcher compares the response body with the expected body via simple string checking. In case of Binary response, byte-by-byte comparison is performed.
func RawBodyMatcher(test *testing.T, expectedBody string, responseBody io.ReadCloser) {
	responseBytes, responseReadErr := io.ReadAll(responseBody)
	if responseReadErr != nil {
		setTestResponseError(test, responseReadErr)
	}
	response := string(responseBytes)

	if !strings.Contains(response, expectedBody) {
		setTestingError(test, responseBody, expectedBody)
	}
}

// NativeBodyMatcher compares the response body as a primitive type(int, int64, float64, bool & time.Time) using a simple equality test.
// Response must match exactly except in case of arrays where array ordering and strictness can be controlled via other options.
func NativeBodyMatcher(test *testing.T, expectedBody string, responseBody io.ReadCloser, isArray, checkArrayCount bool) {
	
	responseBytes, responseReadErr := io.ReadAll(responseBody)
	if responseReadErr != nil {
		setTestResponseError(test, responseReadErr)
	}
	expectedBytes := []byte(expectedBody)
	
	if (isArray){
		matchNativeArray(test, expectedBytes, responseBytes, checkArrayCount)
		return
	}
	if !reflect.DeepEqual(responseBytes, expectedBytes) {
		setTestingError(test, string(responseBytes), string(expectedBytes))
	}
}

func matchNativeArray(test *testing.T, expectedBytes, responseBytes []byte, checkArrayCount bool) {
	var expected, response []any
	expectedErr := json.Unmarshal(expectedBytes, &expected)
	responseErr := json.Unmarshal(responseBytes, &response)
	if expectedErr != nil || responseErr != nil {
		test.Error("error while unmarshalling for comparison")
	}
	if (!checkArrayCount){
		matchNativeArrayValues(test, response, expected)
		return
	}
	if !reflect.DeepEqual(response, expected) {
		setTestingError(test, response, expected)
	}	
}

func matchNativeArrayValues(test *testing.T, response, expected []any) {
	containsFunc := func(slice []any, val any) bool {
		for _, v := range slice {
			if reflect.DeepEqual(v, val) {
				return true
			}
		}
		return false
	}
	for _, v := range response {
		if !containsFunc(expected, v) {
			setTestingError(test, response, expected)
			break
		}
	}
}

// KeysBodyMatcher Checks whether the response body contains the same keys as those specified in the expected body.
// The keys provided can be a subset of the response being received. If any key is absent in the response body, the test fails.
// The test generated will perform deep checking which means if the response object contains nested objects, their keys will also be tested.
func KeysBodyMatcher(test *testing.T, expectedBody string, responseBody io.ReadCloser, checkArrayCount, checkArrayOrder bool) {
	responseBytes, responseErr := io.ReadAll(responseBody)
	if responseErr != nil {
		setTestResponseError(test, responseErr)
	}
	expectedBytes := []byte(expectedBody)

	if !matchKeysAndValues(responseBytes, expectedBytes, checkArrayCount, checkArrayOrder, false) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

// KeysAndValuesBodyMatcher Checks whether the response body contains the same keys and values as those specified in the expected body.
// The keys and values provided can be a subset of the response being received. If any key or value is absent in the response body, the test fails.
// The test generated will perform deep checking which means if the response object contains nested objects, their keys and values will also be tested.
// In case of nested arrays, their ordering and strictness depends on the provided options.
func KeysAndValuesBodyMatcher(test *testing.T, expectedBody string, responseBody io.ReadCloser, checkArrayCount, checkArrayOrder bool) {
	
	responseBytes, responseErr := io.ReadAll(responseBody)
	if responseErr != nil {
		setTestResponseError(test, responseErr)
	}
	expectedBytes := []byte(expectedBody)

	if !matchKeysAndValues(responseBytes, expectedBytes, checkArrayCount, checkArrayOrder, true) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

// matchKeysAndValues is a helper function used by KeysBodyMatcher and KeysAndValuesBodyMatcher
// to compare the bytes for keys and values.
func matchKeysAndValues(actualBytes, expectedBytes []byte, checkArrayCount, checkArrayOrder, checkValues bool) bool {
	actual, expected, ok := extractBytesAsMaps(actualBytes, expectedBytes)
	if !ok {
		return false
	}
	return matchKeysAndValuesAsMap(actual, expected, checkArrayCount, checkArrayOrder, checkValues)
}

// extractBytesAsMaps converts the bytes into maps, if the bytes yield as array,
// then the array indexes will be used as map keys.
func extractBytesAsMaps(actualBytes, expectedBytes []byte) (map[string]any, map[string]any, bool) {
	var actual, expected = make(map[string]any), make(map[string]any)
	var actualErr, expectedErr error

	if len(expectedBytes) > 0 && expectedBytes[0] == '[' {
		var actualArray, expectedArray []any
		actualErr = json.Unmarshal(actualBytes, &actualArray)
		expectedErr = json.Unmarshal(expectedBytes, &expectedArray)
		convertToMap(actualArray, &actual)
		convertToMap(expectedArray, &expected)
	} else {
		actualErr = json.Unmarshal(actualBytes, &actual)
		expectedErr = json.Unmarshal(expectedBytes, &expected)
	}

	if actualErr != nil || expectedErr != nil {
		return nil, nil, false
	}
	return actual, expected, true
}

// matchKeysAndValuesAsMap is a helper function used by matchKeysAndValues
// to compare the JSON keys and values recursively.
func matchKeysAndValuesAsMap(actual, expected map[string]any, checkArrayCount, checkArrayOrder, checkValues bool) bool {
	if checkArrayCount && len(expected) != len(actual) {
		return false
	}
	for key, expectedValue := range expected {
		actualValue := actual[key]
		actualValueKind := reflect.ValueOf(actualValue).Kind()
		if actualValueKind == reflect.Map || actualValueKind == reflect.Array || actualValueKind == reflect.Slice {
			if actualValueKind != reflect.ValueOf(expectedValue).Kind() {
				return false
			}
			actualSubMap, expectedSubMap := extractAnyAsMaps(actualValueKind, actualValue, expectedValue)
			if !matchKeysAndValuesAsMap(actualSubMap, expectedSubMap, checkArrayCount, checkArrayOrder, checkValues) {
				return false
			}
		} else if checkValues && !reflect.DeepEqual(actualValue, expectedValue) {
			return false
		}
	}
	return true
}

// extractAnyAsMaps converts the "any" type values to "map" type values.
// It works only if both types are same and convertible to array or map types
func extractAnyAsMaps(actualValueKind reflect.Kind, actualValue any, expectedValue any) (map[string]any, map[string]any) {
	actualSubMap, expectedSubMap := make(map[string]any), make(map[string]any)
	if actualValueKind != reflect.Map {
		convertToMap(actualValue.([]any), &actualSubMap)
		convertToMap(expectedValue.([]any), &expectedSubMap)
	} else {
		actualSubMap = actualValue.(map[string]any)
		expectedSubMap = expectedValue.(map[string]any)
	}
	return actualSubMap, expectedSubMap
}

// convertToMap add the elements from provided array into the provided map,
// the array indexes will be used as map keys
func convertToMap(array []any, mapRef *map[string]any) {
	if array == nil {
		return
	}
	for i, v := range array {
		(*mapRef)[fmt.Sprint(i)] = v
	}
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
