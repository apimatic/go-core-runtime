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
	json.Unmarshal([]byte(expectedBody), &expected)
	json.Unmarshal([]byte(responseBody), &response)

	if !reflect.DeepEqual(response, expected) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}

func KeysBodyMatcher(test *testing.T, responseBody, expectedBody string) {
	if !compareKeysBody([]byte(responseBody), []byte(expectedBody)) {
		test.Errorf("got \n%v \nbut expected \n%v", responseBody, expectedBody)
	}
}
func compareKeysBody(responseBytes, expectedBytes []byte) bool {
	var expected, response map[string]interface{}
	expectedErr := json.Unmarshal(expectedBytes, &expected)
	responseErr := json.Unmarshal(responseBytes, &response)

	if expectedErr != nil && responseErr != nil{
		return false
	}
	for key,value := range expected {	
		responseValue := response[key]
		if responseValue == nil {
			return false
		}
		switch responseValue.(type) {
		case string : case float32 : case float64 : case bool :
		case int : case int8 : case int16 : case int32 : case int64:
			continue
		default:
			responseBytes,_ := json.Marshal(responseValue)
			expectedBytes,_ := json.Marshal(value)
			if !compareKeysBody(responseBytes, expectedBytes) {
				return false
			}
		}
	}
	return true
}

func KeysAndValuesBodyMatcher(test *testing.T, responseBody, expectedBody string) {
	var expected, response map[string]interface{}
	json.Unmarshal([]byte(expectedBody), &expected)
	json.Unmarshal([]byte(responseBody), &response)

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