package https

import (
	"math"
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

func TestFormEncodeMapNilMap(t *testing.T) {
	result, _ := formEncodeMap(FormParam{"param", "value", nil}, nil)

	expected := []FormParam{
		{"param", "value", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapNilValue(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result, _ := formEncodeMap(FormParam{"param2", nil, nil}, &formParams)

	expected := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMap(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result, _ := formEncodeMap(FormParam{"param2", "value2", nil}, &formParams)

	expected := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2", "value2", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructType(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result, _ := formEncodeMap(FormParam{"param2", GetStruct(), nil}, &formParams)

	expected := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2[Name]", "Bisma", nil},
		{"param2[Employed]", "true", nil},
	}

	var pass bool = true
	for _, res := range result {
		if reflect.DeepEqual(res, FormParam{"param", "value", nil}) {
			continue
		} else if reflect.DeepEqual(res, FormParam{"param1", "value1", nil}) {
			continue
		} else if reflect.DeepEqual(res, FormParam{"param2[Name]", "Bisma", nil}) {
			continue
		} else if reflect.DeepEqual(res, FormParam{"param2[Employed]", "true", nil}) {
			continue
		}
		pass = false
	}
	if !pass {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapMapType(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result, _ := formEncodeMap(FormParam{"param2", map[string]interface{}{"Name": "Bisma"}, nil}, &formParams)

	expected := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2[Name]", "Bisma", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapSliceType(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result, _ := formEncodeMap(FormParam{"param2", []string{"Name", "Bisma"}, nil}, &formParams)

	expected := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2[0]", "Name", nil},
		{"param2[1]", "Bisma", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapInterfaceSliceType(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result, _ := formEncodeMap(FormParam{"param2", []interface{}{"Name", "Bisma"}, nil}, &formParams)

	expected := []FormParam{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2[0]", "Name", nil},
		{"param2[1]", "Bisma", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructTypeError(t *testing.T) {
	formParams := []FormParam{
		{"param", "value", nil},
	}
	ptr := math.Inf(1)
	_, err := formEncodeMap(FormParam{"param2", &ptr, nil}, &formParams)

	if err == nil {
		t.Errorf("The code should get error because input cannot be converted to struct.")
	}
}

func TestPrepareFormFieldsNil(t *testing.T) {
	result, _ := prepareFormFields(FormParam{"param", "value", nil}, nil)

	expected := url.Values{}
	expected.Add("param", "value")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFields(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", "value", nil}, input)

	expected := url.Values{}
	expected.Add("param", "val")
	expected.Add("param", "val1")
	expected.Add("param2", "value")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsStringSlice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []string{"value", "value1"}, nil}, input)

	expected := input
	expected.Add("param2", "value")
	expected.Add("param2", "value1")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsIntSlice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []int{1, 2}, nil}, input)

	expected := input
	expected.Add("param2", "1")
	expected.Add("param2", "2")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsInt16Slice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []int16{1, 2}, nil}, input)

	expected := input
	expected.Add("param2", "1")
	expected.Add("param2", "2")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsInt32Slice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []int32{1, 2}, nil}, input)

	expected := input
	expected.Add("param2", "1")
	expected.Add("param2", "2")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsInt64Slice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []int64{1, 2}, nil}, input)

	expected := input
	expected.Add("param2", "1")
	expected.Add("param2", "2")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsBoolSlice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []bool{false, true}, nil}, input)

	expected := input
	expected.Add("param2", "false")
	expected.Add("param2", "true")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsFloat32Slice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []float32{1.2, 2.1}, nil}, input)

	expected := input
	expected.Add("param2", "1.2")
	expected.Add("param2", "2.1")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsFloat64Slice(t *testing.T) {
	input := url.Values{}
	input.Add("param", "val")
	input.Add("param", "val1")
	result, _ := prepareFormFields(FormParam{"param2", []float64{1.1111, 2.1111}, nil}, input)

	expected := input
	expected.Add("param2", "1.1111")
	expected.Add("param2", "2.1111")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsFloat64Pointer(t *testing.T) {
	floatV := math.Inf(1)
	result, err := prepareFormFields(FormParam{"param", &(floatV), nil}, nil)

	if err == nil {
		t.Errorf("Failed:\nExpected: nil \nGot: %v", result)
	}
}

func TestPrepareMultipartFields(t *testing.T) {
	bytes, str, _ := prepareMultipartFields([]FormParam{{"param", "value", nil}})

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithPointer(t *testing.T) {
	floatV := math.Inf(0)
	bytes, str, _ := prepareMultipartFields([]FormParam{{"param", &floatV, nil}})

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFile(t *testing.T) {
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}
	bytes, _, _ := prepareMultipartFields([]FormParam{{"param", file, nil}})

	if !strings.Contains(bytes.String(), `filename=googles-new-logo`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFileError(t *testing.T) {
	bytes, _, _ := prepareMultipartFields([]FormParam{{"param", nil, nil}})

	if !strings.Contains(bytes.String(), `null`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}
