package https

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/internal"
	"github.com/google/uuid"
)

type Person struct {
	Name     string
	Employed bool
}

func GetStruct() Person {
	return Person{Name: "Bisma", Employed: true}
}

func TestStructToMap(t *testing.T) {
	result, _ := structToAny(GetStruct())

	expected := map[string]any{
		"Name":     "Bisma",
		"Employed": true,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestStructToMapMarshallingError(t *testing.T) {
	result, err := structToAny(math.Inf(1))

	if err == nil && result != nil {
		t.Error("Failed:\nExpected error in marshalling infinity number")
	}
}

func TestToMapNilMap(t *testing.T) {
	param := formParam{"param", "value", nil, Indexed}
	result, _ := param.toMap()

	expected := map[string][]string{
		"param": {"value"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapNilValue(t *testing.T) {
	param := formParam{"param", nil, nil, Indexed}
	result, _ := param.toMap()

	expected := make(map[string][]string)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeUUIDValue(t *testing.T) {
	uuidVal, _ := uuid.ParseBytes([]byte("992bf4b9-c900-4850-9992-107b2f9df928"))
	param := formParam{"uuid-param", uuidVal, nil, Indexed}
	result, _ := param.toMap()

	expected := make(map[string][]string)
	expected[param.key] = []string{"992bf4b9-c900-4850-9992-107b2f9df928"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestFormEncodeMapStructType(t *testing.T) {
	param := formParam{"param2", GetStruct(), nil, Indexed}
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
	params := formParams{{"param", "value", nil, Indexed}}
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
	params := formParams{{"param2", "value", nil, Indexed}}
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
	params := formParams{{"param2", []bool{false, true}, nil, Indexed}}
	_ = params.prepareFormFields(result)

	expected := result
	expected.Add("param2", "false")
	expected.Add("param2", "true")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsAnySlice(t *testing.T) {
	anySlice := []any{
		any("Item1"),
		any("Item2"),
	}

	params := formParams{
		{"anySlice", anySlice, nil, Csv},
	}
	result := url.Values{}
	_ = params.prepareFormFields(result)

	expected := url.Values{}
	expected.Add("anySlice", fmt.Sprintf("%v,%v", anySlice[0], anySlice[1]))

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsEnumSlice(t *testing.T) {
	stringEnums := []internal.MonthNameEnum{
		internal.MonthNameEnum_JANUARY,
		internal.MonthNameEnum_FEBRUARY,
		internal.MonthNameEnum_MARCH,
	}
	numberEnums := []internal.MonthNumberEnum{
		internal.MonthNumberEnum_JANUARY,
		internal.MonthNumberEnum_FEBRUARY,
		internal.MonthNumberEnum_MARCH,
	}

	params := formParams{
		{"stringEnums", stringEnums, nil, Csv},
		{"numberEnums", numberEnums, nil, Csv},
	}
	result := url.Values{}
	_ = params.prepareFormFields(result)

	expected := url.Values{}
	expected.Add("stringEnums", fmt.Sprintf("%v,%v,%v", stringEnums[0], stringEnums[1], stringEnums[2]))
	expected.Add("numberEnums", fmt.Sprintf("%v,%v,%v", numberEnums[0], numberEnums[1], numberEnums[2]))

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareFormFieldsFloat64Pointer(t *testing.T) {
	floatV := math.Inf(1)
	params := formParams{{"param", &floatV, nil, Indexed}}
	result := url.Values{}
	err := params.prepareFormFields(result)
	if err == nil {
		t.Errorf("Failed:\nExpected: nil \nGot: %v", result)
	}
}

func TestPrepareMultipartFieldsString(t *testing.T) {
	header := http.Header{}
	header.Add(CONTENT_TYPE_HEADER, TEXT_CONTENT_TYPE)
	params := formParams{{"param", "value", header, Indexed}}
	bytes, str, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, MULTIPART_CONTENT_TYPE) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFields(t *testing.T) {
	header := http.Header{}
	header.Add(CONTENT_TYPE_HEADER, TEXT_CONTENT_TYPE)
	params := formParams{{"param", 40, header, Indexed}}
	bytes, str, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, MULTIPART_CONTENT_TYPE) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithPointer(t *testing.T) {
	floatV := math.Inf(0)
	params := formParams{{"param", &floatV, nil, Indexed}}
	bytes, str, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `name="param"`) && !strings.Contains(str, MULTIPART_CONTENT_TYPE) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}

func TestPrepareMultipartFieldsWithFile(t *testing.T) {
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}
	header := http.Header{}
	header.Add(CONTENT_TYPE_HEADER, "image/png")
	params := formParams{{"param", file, header, Indexed}}
	bytes, _, _ := params.prepareMultipartFields()

	if !strings.Contains(bytes.String(), `filename=googles-new-logo`) {
		t.Errorf("Failed:\nGot: %v", bytes.String())
	}
}
