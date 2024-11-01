package testHelper_test

import (
	"bytes"
	"encoding/json"
	"github.com/apimatic/go-core-runtime/testHelper"
	"io"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

func getValueReader(val any) io.ReadCloser {
	byt, _ := json.Marshal(val)
	return io.NopCloser(bytes.NewBuffer(byt))
}

// Native Body Matcher Tests
func TestNativeBodyMatcherNumber(t *testing.T) {
	expected := `4`
	var result = 4
	testHelper.NativeBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestNativeBodyMatcherPrecision(t *testing.T) {
	expected := `4.11`
	var result float32 = 4.11
	testHelper.NativeBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestNativeBodyMatcherLong(t *testing.T) {
	expected := `411111111111111111`
	var result int64 = 411111111111111111
	testHelper.NativeBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestNativeBodyMatcherBoolean(t *testing.T) {
	expected := `true`
	testHelper.NativeBodyMatcher(t, expected, getValueReader(true), false, false)
}

func TestNativeBodyMatcherStringSlice(t *testing.T) {
	expected := `["Tuesday", "Saturday", "Wednesday", "Monday", "Sunday"]`
	var result = []string{
		"Tuesday", "Saturday", "Wednesday", "Monday", "Sunday",
	}
	testHelper.NativeBodyMatcher(t, expected, getValueReader(result), true, true)
}

func TestNativeBodyMatcherIntSlice(t *testing.T) {
	expected := `[1,2,3,4,5]`
	var result = []int{
		1, 2, 3, 4, 5,
	}
	testHelper.NativeBodyMatcher(t, expected, getValueReader(result), true, false)
}

func TestNativeBodyMatcherBooleanError(t *testing.T) {
	expected := `nil`
	testHelper.NativeBodyMatcher(&testing.T{}, expected, getValueReader(true), false, false)
}

// Raw Body Matcher Tests
func TestRawBodyMatcherIntSlice(t *testing.T) {
	expected := `[1,2,3,4,5]`
	var result = []int{
		1, 2, 3, 4, 5,
	}
	testHelper.RawBodyMatcher(t, expected, getValueReader(result))
}

func TestRawBodyMatcherBooleanError(t *testing.T) {
	expected := `nil`
	testHelper.RawBodyMatcher(&testing.T{}, expected, getValueReader(true))
}

// Is Same File Tests
func TestIsSameFile(t *testing.T) {
	expected := `https://www.gstatic.com/webp/gallery/1.jpg`
	result, err := https.GetFile(expected)
	if err != nil {
		t.Error("Error fetching File from ", expected)
	}
	testHelper.IsSameAsFile(t, expected, result.File)
}

func TestIsSameFileError(t *testing.T) {
	expected := `http://localhost:3000/response/image`
	result, _ := https.GetFile("https://play.google.com/store/apps/dev?id=5700313618786177705&hl=en_US&gl=US")
	testHelper.IsSameAsFile(&testing.T{}, expected, result.File)
}

func TestIsSameFileErrorURL(t *testing.T) {
	expected := `http://response/image`
	result, _ := https.GetFile(`http://localhost:3000/response/image`)
	testHelper.IsSameAsFile(&testing.T{}, expected, result.File)
}

// Slice to Comma Separated String Tests
func TestSliceToCommaSeparatedString(t *testing.T) {
	expected := `{"isMap": false,"id": "5a9fcb01caacc310dc6bab50"}`
	testHelper.SliceToCommaSeparatedString(expected)
}

// Keys And Values Body Matcher Tests
type Response struct {
	IsMap           bool         `json:"isMap"`
	Attributes      Attributes   `json:"attributes"`
	AttributesArray []Attributes `json:"attributesArray"`
	Id              string       `json:"id"`
}
type Attributes struct {
	Id string `json:"id"`
}

func TestKeysAndValuesBodyMatcherEmpty(t *testing.T) {
	expected := `{}`
	testHelper.KeysAndValuesBodyMatcher(t, expected, getValueReader(nil), false, false)
}

func TestKeysAndValuesBodyMatcherEmptyArray(t *testing.T) {
	expected := `[]`
	testHelper.KeysAndValuesBodyMatcher(t, expected, getValueReader(nil), false, false)
}

func TestKeysAndValuesBodyMatcherArray(t *testing.T) {
	expected := `["some string", 123]`
	result := []any{"some string", 123}
	testHelper.KeysAndValuesBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysAndValuesBodyMatcherObject(t *testing.T) {
	expected := `{"id": "5a9fcb01caacc310dc6bab51"}`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	testHelper.KeysAndValuesBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysAndValuesBodyMatcherNestedObject(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributes": {
          "id": "5a9fcb01caacc310dc6bab51"
        },
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		Attributes: Attributes{
			Id: "5a9fcb01caacc310dc6bab51",
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysAndValuesBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysAndValuesBodyMatcherNestedArray(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributesArray": [
			{
          		"id": "5a9fcb01caacc310dc6bab51"
			}
		],
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		AttributesArray: []Attributes{
			{Id: "5a9fcb01caacc310dc6bab51"},
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysAndValuesBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysAndValuesBodyMatcherNestedObjectValueError(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributes": {
          "id": "5a9fcb01caacc310dc6bab51"
        },
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		Attributes: Attributes{
			Id: "5a9fcb01caacc0dc6bab51",
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysAndValuesBodyMatcher(&testing.T{}, expected, getValueReader(result), false, false)
}

func TestKeysAndValuesBodyMatcherNestedObjectTypeError(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributes": "5a9fcb01caacc310dc6bab51",
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		Attributes: Attributes{
			Id: "5a9fcb01caacc0dc6bab51",
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysAndValuesBodyMatcher(&testing.T{}, expected, getValueReader(result), false, false)
}

func TestKeysAndValuesBodyMatcherNestedObjectArrayCountError(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributes": {
        },
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		Attributes: Attributes{
			Id: "5a9fcb01caacc310dc6bab51",
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysAndValuesBodyMatcher(&testing.T{}, expected, getValueReader(result), true, false)
}

func TestKeysAndValuesBodyMatcherUnmarshallingError(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributes": {
        },
        "id": "5a9fcb01caacc310dc6bab50"
    `
	result := Response{
		IsMap: false,
		Attributes: Attributes{
			Id: "5a9fcb01caacc310dc6bab51",
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysAndValuesBodyMatcher(&testing.T{}, expected, getValueReader(result), true, false)
}

// Keys Body Matcher Tests
func TestKeysBodyMatcherEmpty(t *testing.T) {
	expected := `{}`
	testHelper.KeysBodyMatcher(t, expected, getValueReader(nil), false, false)
}

func TestKeysBodyMatcherEmptyArray(t *testing.T) {
	expected := `[]`
	testHelper.KeysBodyMatcher(t, expected, getValueReader(nil), false, false)
}

func TestKeysBodyMatcherArray(t *testing.T) {
	expected := `["some string", 123]`
	result := []any{"123", 765}
	testHelper.KeysBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysBodyMatcherObject(t *testing.T) {
	expected := `{"id": "5a9fcb01caacc310dc6bab51"}`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	testHelper.KeysBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysBodyMatcherNestedArray(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributesArray": [
			{
          		"id": "5a9fcb01caacc310dc6bab51"
			}
		],
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		AttributesArray: []Attributes{
			{Id: "5a9fcb01caacc310dc6bab51"},
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysBodyMatcherNestedObject(t *testing.T) {
	expected := `{
        "isMap": false,
        "attributes": {
          "id": "5a9fcb01caacc310dc6bab51"
        },
        "id": "5a9fcb01caacc310dc6bab50"
    }`
	result := Response{
		IsMap: false,
		Attributes: Attributes{
			Id: "5a9fcb01caacc310dc6bab51",
		},
		Id: "5a9fcb01caacc310dc6bab50",
	}
	testHelper.KeysBodyMatcher(t, expected, getValueReader(result), false, false)
}

func TestKeysBodyMatcherObjectError(t *testing.T) {
	expected := `{"idd": "nil"}`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	testHelper.KeysBodyMatcher(&testing.T{}, expected, getValueReader(result), false, false)
}

func TestKeysBodyMatcherUnmarshallingError(t *testing.T) {
	expected := `{"idd": "nil"`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	testHelper.KeysBodyMatcher(&testing.T{}, expected, getValueReader(result), false, false)
}
