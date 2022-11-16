package testHelper

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

func NativeBodyMatcher(test *testing.T, expectedBody string, responseObject any) {
	responseBytes,_ := json.Marshal(responseObject)
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

	responseBytes,_ := json.Marshal(responseObject)
	var response, expected map[string]interface{}
	responseErr := json.Unmarshal(responseBytes, &response)
	expectedErr := json.Unmarshal([]byte(expectedBody), &expected)

	if responseErr != nil && expectedErr != nil{
		test.Error("Error while Unmarshalling")
	}

	if !matchKeysAndValuesBody(response, expected, checkArrayCount, checkArrayOrder, false) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

func KeysAndValuesBodyMatcher(test *testing.T, expectedBody string, responseObject any, checkArrayCount, checkArrayOrder bool) {
	responseBytes,_ := json.Marshal(responseObject)
	var response, expected map[string]interface{}
	responseErr := json.Unmarshal(responseBytes, &response)
	expectedErr := json.Unmarshal([]byte(expectedBody), &expected)

	if responseErr != nil && expectedErr != nil{
		test.Error("Error while Unmarshalling")
	}

	if !matchKeysAndValuesBody(response, expected, checkArrayCount, checkArrayOrder, true) {
		test.Errorf("got \n%v \nbut expected \n%v", string(responseBytes), expectedBody)
	}
}

func matchKeysAndValuesBody(response, expected map[string]interface{}, checkArrayCount, checkArrayOrder, checkValues bool) bool {
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
			if !matchKeysAndValuesBody(responseSubMap, expectedSubMap, checkArrayCount, checkArrayOrder, checkValues) {
				return false
			}
		} else if checkValues && !reflect.DeepEqual(responseValue, value) {
			return false
		}
	}
	return true
}

func RawBodyMatcher(test *testing.T, expectedBody string, responseObject any) {
	responseBytes,_ := json.Marshal(responseObject)
	responseBody := string(responseBytes)

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
