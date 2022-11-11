package testHelper

import (
	"encoding/json"
	"net/http"
	"testing"
)

type TestHeader struct {
	CheckValue  bool            `json:"CheckValue"`     
    Name  		string 			`json:"Name"`     
    Value     	string         `json:"Value"`        
}

func CheckResponseHeaders(t *testing.T, respHeaders http.Header, expectedHeaders string, allowExtra bool) {
	
	var expectedHeadersMaps []TestHeader
	json.Unmarshal([]byte(expectedHeaders), &expectedHeadersMaps)
	expectedHeadersCount := 0

	for index, expectedHeadersMap := range expectedHeadersMaps {
		expectedHeadersCount += index

		respValue := respHeaders.Get(expectedHeadersMap.Name)
		if(respValue == "") {
			t.Errorf("expected header '%v' does not exists in response", expectedHeadersMap.Name)
			break
		} else if expectedHeadersMap.CheckValue && respValue != expectedHeadersMap.Value {
			t.Errorf("response does not contains same value of expected header '%v'", expectedHeadersMap.Name)
			break
		}
    }
	if !allowExtra && len(respHeaders) != expectedHeadersCount {
		t.Errorf("response contains other headers than those listed in the expected headers list")
	}
}