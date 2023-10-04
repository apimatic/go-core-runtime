package https

import (
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type Person struct {
	Name     string
	Employed bool
}

func GetStruct() Person {
	return Person{Name: "Bisma", Employed: true}
}

func TestStructToMap(t *testing.T) {
	result, _ := structToMap(GetStruct())

	expected := map[string]interface{}{
		"Name":     "Bisma",
		"Employed": true,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestStructToMapMarshallingError(t *testing.T) {
	result, err := structToMap(math.Inf(1))

	if err == nil && result != nil {
		t.Error("Failed:\nExpected error in marshalling infinity number")
	}
}

func TestToMapNilMap(t *testing.T) {
	param := FormParam{"param", "value", nil}
	result, _ := param.toMap()

	expected := map[string]string{
		"param": "value",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapNilValue(t *testing.T) {
	param := FormParam{"param", nil, nil}
	result, _ := param.toMap()

	expected := make(map[string]string)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructType(t *testing.T) {
	param := FormParam{"param2", GetStruct(), nil}
	result, _ := param.toMap()

	expected := FormParams{
		{"param2[Name]", "Bisma", nil},
		{"param2[Employed]", "true", nil},
	}

	if len(result) != len(expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsNil(t *testing.T) {
	params := FormParams{{"param", "value", nil}}
	result := url.Values{}
	_ = params.prepareFormFields(result)

	expected := url.Values{}
	expected.Add("param", "value")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFields(t *testing.T) {
	result := url.Values{}
	result.Add("param", "val")
	result.Add("param", "val1")
	params := FormParams{{"param2", "value", nil}}
	_ = params.prepareFormFields(result)

	expected := url.Values{}
	expected.Add("param", "val")
	expected.Add("param", "val1")
	expected.Add("param2", "value")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsBoolSlice(t *testing.T) {
	result := url.Values{}
	result.Add("param", "val")
	result.Add("param", "val1")
	params := FormParams{{"param2", []bool{false, true}, nil}}
	_ = params.prepareFormFields(result)

	expected := result
	expected.Add("param2", "false")
	expected.Add("param2", "true")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsFloat64Pointer(t *testing.T) {
	floatV := math.Inf(1)
	params := FormParams{{"param", &floatV, nil}}
	result := url.Values{}
	err := params.prepareFormFields(result)
	if err == nil {
		t.Errorf("Failed:\nExpected: nil \nGot: %v", result)
	}
}

func TestPrepareMultipartFieldsString(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", TEXT_CONTENT_TYPE)
	params := FormParams{{"param", "value", header}}
	bytes, str, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFields(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", TEXT_CONTENT_TYPE)
	params := FormParams{{"param", 40, header}}
	bytes, str, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithPointer(t *testing.T) {
	floatV := math.Inf(0)
	params := FormParams{{"param", &floatV, nil}}
	bytes, str, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFile(t *testing.T) {
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}
	header := http.Header{}
	header.Add("Content-Type", "image/png")
	params := FormParams{{"param", file, header}}
	bytes, _, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `filename=googles-new-logo`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}
