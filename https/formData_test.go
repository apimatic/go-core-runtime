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

func TestFormEncodeMapNilMap(t *testing.T) {

	result := FormParams{
		FormParam{"param", "value", nil},
	}

	expected := FormParams{
		{"param", "value", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapNilValue(t *testing.T) {
	result := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result.Add(FormParam{"param2", nil, nil})

	expected := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMap(t *testing.T) {
	result := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result.Add(FormParam{"param2", "value2", nil})

	expected := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2", "value2", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructType(t *testing.T) {
	result := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result.Add(FormParam{"param2", GetStruct(), nil})

	expected := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2[Name]", "Bisma", nil},
		{"param2[Employed]", "true", nil},
	}

	if len(result) != len(expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapMapType(t *testing.T) {
	result := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result.Add(FormParam{"param2", map[string]interface{}{"Name": "Bisma"}, nil})

	expected := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
		{"param2", "map[Name:Bisma]", nil},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapSliceType(t *testing.T) {
	result := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result.Add(FormParam{"param2", []string{"Name", "Bisma"}, nil})

	expected := FormParams{
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
	result := FormParams{
		{"param", "value", nil},
		{"param1", "value1", nil},
	}
	result.Add(FormParam{"param2", []interface{}{"Name", "Bisma"}, nil})

	expected := FormParams{
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
	result := FormParams{
		{"param", "value", nil},
	}
	ptr := math.Inf(1)
	result.Add(FormParam{"param2", &ptr, nil})

	// if err == nil {
	// 	t.Errorf("The code should get error because input cannot be converted to struct.")
	// }
}

func TestPrepareFormFieldsNil(t *testing.T) {
	formParams := FormParams{
		FormParam{"param", "value", nil},
	}
	result, _ := formParams.prepareFormFields(nil)

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
	formParams := FormParams{
		FormParam{"param2", "value", nil},
	}
	result, _ := formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []string{"value", "value1"}, nil},
	}
	result, _ := formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []int{1, 2}, nil},
	}
	result, _ := formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []int16{1, 2}, nil},
	}
	result, _ := formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []int32{1, 2}, nil},
	}
	result, _ := formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []int64{1, 2}, nil},
	}
	result, _ :=formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []bool{false, true}, nil},
	}
	result, _ :=formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []float32{1.2, 2.1}, nil}, 
	}
	result, _ :=formParams.prepareFormFields(input)

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
	formParams := FormParams{
		FormParam{"param2", []float64{1.1111, 2.1111}, nil}, 
	}
	result, _ := formParams.prepareFormFields(input)

	expected := input
	expected.Add("param2", "1.1111")
	expected.Add("param2", "2.1111")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsFloat64Pointer(t *testing.T) {
	floatV := math.Inf(1)
	formParams := FormParams{
		FormParam{"param", &(floatV), nil}, 
	}
	result, err := formParams.prepareFormFields(nil)

	if err == nil {
		t.Errorf("Failed:\nExpected: nil \nGot: %v", result)
	}
}

func TestPrepareMultipartFieldsString(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", TEXT_CONTENT_TYPE)
	formParams := FormParams{{"param", "value", header}}
	bytes, str, _ := formParams.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFields(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", TEXT_CONTENT_TYPE)
	formParams := FormParams{{"param", 40, header}}
	bytes, str, _ := formParams.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithPointer(t *testing.T) {
	floatV := math.Inf(0)
	formParams := FormParams{{"param", &floatV, nil}}
	bytes, str, _ := formParams.prepareMultipartFields()

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
	formParams := FormParams{{"param", file, header}}
	bytes, _, _ := formParams.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `filename=googles-new-logo`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFileError(t *testing.T) {
	formParams := FormParams{{"param", nil, nil}}
	bytes, _, _ := formParams.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `null`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}
