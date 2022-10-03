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
	result := formEncodeMap("param", "value", nil)

	expected := []map[string]interface{}{
		{"param": "value"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapNilValue(t *testing.T) {
	mapInput := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
	}
	result := formEncodeMap("param2", nil, &mapInput)

	expected := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
		{},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMap(t *testing.T) {
	mapInput := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
	}
	result := formEncodeMap("param2", "value2", &mapInput)

	expected := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
		{"param2": "value2"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructType(t *testing.T) {
	mapInput := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
	}
	result := formEncodeMap("param2", GetStruct(), &mapInput)

	expected := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
		{"param2[Name]": "Bisma"},
		{"param2[Employed]": "true"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapMapType(t *testing.T) {
	mapInput := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
	}
	result := formEncodeMap("param2", map[string]interface{}{"Name": "Bisma"}, &mapInput)

	expected := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
		{"param2[Name]": "Bisma"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapSliceType(t *testing.T) {
	mapInput := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
	}
	result := formEncodeMap("param2", []string{"Name", "Bisma"}, &mapInput)

	expected := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
		{"param2[0]": "Name"},
		{"param2[1]": "Bisma"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapInterfaceSliceType(t *testing.T) {
	mapInput := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
	}
	result := formEncodeMap("param2", []interface{}{"Name", "Bisma"}, &mapInput)

	expected := []map[string]interface{}{
		{"param": "value", "param1": "value1"},
		{"param2[0]": "Name"},
		{"param2[1]": "Bisma"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructTypeError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should panic because input cannot be converted to struct.")
		}
	}()

	mapInput := []map[string]interface{}{
		{"param": "value"},
	}
	ptr := math.Inf(1)
	formEncodeMap("param2", &ptr, &mapInput)
}

func TestPrepareFormFieldsNil(t *testing.T) {
	result := PrepareFormFields("param", "value", nil)

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
	result := PrepareFormFields("param2", "value", input)

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
	result := PrepareFormFields("param2", []string{"value", "value1"}, input)

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
	result := PrepareFormFields("param2", []int{1, 2}, input)

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
	result := PrepareFormFields("param2", []int16{1, 2}, input)

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
	result := PrepareFormFields("param2", []int32{1, 2}, input)

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
	result := PrepareFormFields("param2", []int64{1, 2}, input)

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
	result := PrepareFormFields("param2", []bool{false, true}, input)

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
	result := PrepareFormFields("param2", []float32{1.2, 2.1}, input)

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
	result := PrepareFormFields("param2", []float64{1.1111, 2.1111}, input)

	expected := input
	expected.Add("param2", "1.1111")
	expected.Add("param2", "2.1111")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareMultipartFields(t *testing.T) {
	bytes, str := PrepareMultipartFields(map[string]interface{}{"param": "value"})

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, "multipart/form-data") {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFile(t *testing.T) {
	file := GetFile("https://www.google.com/doodles/googles-new-logo")
	bytes, _ := PrepareMultipartFields(map[string]interface{}{"param": file})

	if !strings.Contains(bytes.String(), `filename="googles-new-logo"`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFileError(t *testing.T) {
	bytes, _ := PrepareMultipartFields(map[string]interface{}{"param": nil})

	if !strings.Contains(bytes.String(), `null`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}
