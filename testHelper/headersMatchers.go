package testHelper

import (
	"net/http"
	"testing"
)

type TestHeader struct {
	CheckValue  bool            `json:"CheckValue"`     
    Name  		string 			`json:"Name"`     
    Value     	string         `json:"Value"`        
}

func NewTestHeader(checkValue bool, name, value string) TestHeader {
	return TestHeader{
        CheckValue: checkValue,
		Name: name,
		Value: value,
    }
}

func CheckResponseHeaders(t *testing.T, responseHeaders http.Header, expectedHeadersList []TestHeader, allowExtraHeaders bool) {
	
	for _, expectedHeader := range expectedHeadersList {
		respValue := responseHeaders.Get(expectedHeader.Name)
		if(respValue == "") {
			t.Errorf("expected header '%v' does not exists in response", expectedHeader.Name)
			break
		} else if expectedHeader.CheckValue && respValue != expectedHeader.Value {
			t.Errorf("response does not contains same value of expected header '%v'", expectedHeader.Name)
			break
		}
    }
	if !allowExtraHeaders && len(responseHeaders) != len(expectedHeadersList) {
		t.Errorf("response contains other headers than those listed in the expected headers list")
	}
}