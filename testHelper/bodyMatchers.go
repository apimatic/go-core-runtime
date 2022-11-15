package testHelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

func NativeBodyMatcher(test *testing.T, responseBody, expectedBody string) {
	var expected, response interface{}
	expectedError := json.Unmarshal([]byte(expectedBody), &expected)
	responseError := json.Unmarshal([]byte(responseBody), &response)

	if expectedError != nil || responseError != nil {
		test.Error("Error while Unmarshalling")
	}

	if !reflect.DeepEqual(response, expected) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}

func KeysBodyMatcher(test *testing.T, responseBody, expectedBody string) {
	var expected, response map[string]interface{}
	expectedError := json.Unmarshal([]byte(expectedBody), &expected)
	responseError := json.Unmarshal([]byte(responseBody), &response)

	if expectedError != nil || responseError != nil {
		test.Error("Error while Unmarshalling")
	}

	responseStringSlice := toStringSlice(getMapKeys(response))
	sort.Strings(responseStringSlice)

	expectedStringSlice := toStringSlice(getMapKeys(expected))
	sort.Strings(expectedStringSlice)

	if !reflect.DeepEqual(responseStringSlice, expectedStringSlice) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}

func KeysAndValuesBodyMatcher(test *testing.T, responseBody, expectedBody string) {
	var expected, response map[string]interface{}
	expectedError := json.Unmarshal([]byte(expectedBody), &expected)
	responseError := json.Unmarshal([]byte(responseBody), &response)

	if expectedError != nil || responseError != nil {
		test.Error("Error while Unmarshalling")
	}

	if !reflect.DeepEqual(response, expected) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}

func RawBodyMatcher(test *testing.T, responseBody, expectedBody string) {
	if responseBody != expectedBody {
		test.Errorf("got \n%v \nbut expected %v", responseBody, expectedBody)
	}
}

func IsSameFile(test *testing.T, expectedFileURL string, responseFile https.FileWrapper) {
	expectedFile := https.GetFile(expectedFileURL).File
	if !reflect.DeepEqual(responseFile, expectedFile) {
		test.Error("Response File does not match the File received")
	}
}

func SliceToCommaSeparatedString(slice interface{}) string {
	return strings.Join(strings.Split(fmt.Sprint(slice), " "), ",")
}

func toStringSlice(input []reflect.Value) []string {
	stringSlice := make([]string, 0)
	for _, v := range input {
		stringSlice = append(stringSlice, v.Interface().(string))
	}
	return stringSlice
}

func getMapKeys(input map[string]interface{}) []reflect.Value {
	keys := reflect.ValueOf(input).MapKeys()
	for _, v := range input {
		x := reflect.ValueOf(v)
		if x.Kind() == reflect.Map {
			keys = append(keys, getMapKeys(v.(map[string]interface{}))...)
		}
	}
	return keys
}
