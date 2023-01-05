package testHelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

func NativeBodyMatcher(test *testing.T, expectedBody string, responseObject any) {
	responseBytes, _ := json.Marshal(responseObject)
	var expected, response interface{}
	expectedError := json.Unmarshal([]byte(expectedBody), &expected)
	responseError := json.Unmarshal(responseBytes, &response)

	if expectedError != nil || responseError != nil {
		test.Error("Error while Unmarshalling")
	}

	if !reflect.DeepEqual(response, expected) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

func KeysBodyMatcher(test *testing.T, expectedBody string, responseObject any, checkArrayCount, checkArrayOrder bool) {
	responseBytes, _ := json.Marshal(responseObject)
	var response, expected map[string]interface{}
	responseErr := json.Unmarshal(responseBytes, &response)
	expectedErr := json.Unmarshal([]byte(expectedBody), &expected)

	if responseErr != nil || expectedErr != nil {
		test.Error("Error while Unmarshalling")
	}

	if !matchKeysAndValues(response, expected, checkArrayCount, checkArrayOrder, false) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

func KeysAndValuesBodyMatcher(test *testing.T, expectedBody string, responseObject any, checkArrayCount, checkArrayOrder bool) {
	responseBytes, _ := json.Marshal(responseObject)
	var response, expected map[string]interface{}
	responseErr := json.Unmarshal(responseBytes, &response)
	expectedErr := json.Unmarshal([]byte(expectedBody), &expected)

	if responseErr != nil || expectedErr != nil {
		test.Error("Error while Unmarshalling")
	}

	if !matchKeysAndValues(response, expected, checkArrayCount, checkArrayOrder, true) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

func matchKeysAndValues(response, expected map[string]interface{}, checkArrayCount, checkArrayOrder, checkValues bool) bool {
	if checkArrayCount && len(expected) != len(response) {
		return false
	}
	for key, value := range expected {
		responseValue := response[key]
		if responseValue == nil {
			return false
		}
		if reflect.ValueOf(responseValue).Kind() == reflect.Map {
			if reflect.ValueOf(value).Kind() != reflect.Map {
				return false
			}
			responseSubMap := responseValue.(map[string]interface{})
			expectedSubMap := value.(map[string]interface{})
			if !matchKeysAndValues(responseSubMap, expectedSubMap, checkArrayCount, checkArrayOrder, checkValues) {
				return false
			}
		} else if checkValues && !reflect.DeepEqual(responseValue, value) {
			return false
		}
	}
	return true
}

func RawBodyMatcher(test *testing.T, expectedBody string, responseObject any) {
	responseBytes, _ := json.Marshal(responseObject)
	responseBody := string(responseBytes)

	if !strings.Contains(responseBody, expectedBody) {
		test.Errorf("got \n%v \nbut expected %v", responseBody, expectedBody)
	}
}

func IsSameAsFile(test *testing.T, expectedFileURL string, responseFileBytes []byte) {
	expectedFile, err := https.GetFile(expectedFileURL)
	if err != nil {
		test.Errorf("Cannot get the file: %v", err)
	}
	IsSameInputBytes(test, expectedFile.File, responseFileBytes)
}

func IsSameInputBytes(test *testing.T, expectedBytes []byte, receivedBytes []byte) {
	if !reflect.DeepEqual(expectedBytes, receivedBytes) {
		test.Error("Recieved bytes donot match the bytes expected")
	}
}

func SliceToCommaSeparatedString(slice interface{}) string {
	return strings.Join(strings.Split(fmt.Sprint(slice), " "), ",")
}
