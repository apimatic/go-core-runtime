package testHelper

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func KeysBodyMatcher(test *testing.T, responseBody, expectedBody string, checkArrayCount, checkArrayOrder bool) {

	var response, expected map[string]interface{}
	responseErr := json.Unmarshal([]byte(responseBody), &response)
	expectedErr := json.Unmarshal([]byte(expectedBody), &expected)

	if responseErr != nil && expectedErr != nil{
		test.Error("Error while Unmarshalling")
	}

	if !compareKeysBody(response, expected, checkArrayCount, checkArrayOrder, false) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}

func KeysAndValuesBodyMatcher(test *testing.T, responseBody, expectedBody string, checkArrayCount, checkArrayOrder bool) {
	var response, expected map[string]interface{}
	responseErr := json.Unmarshal([]byte(responseBody), &response)
	expectedErr := json.Unmarshal([]byte(expectedBody), &expected)

	if responseErr != nil && expectedErr != nil{
		test.Error("Error while Unmarshalling")
	}

	if !compareKeysBody(response, expected, checkArrayCount, checkArrayOrder, true) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}

func compareKeysBody(response, expected map[string]interface{}, checkArrayCount, checkArrayOrder, checkValues bool) bool {
	if checkArrayCount && len(expected) != len(response) {
		return false
	}
	for key,value := range expected {	
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
			if !compareKeysBody(responseSubMap, expectedSubMap, checkArrayCount, checkArrayOrder, checkValues) {
				return false
			}
		} else if checkValues && !reflect.DeepEqual(responseValue, value) {
			return false
		}
	}
	return true
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
