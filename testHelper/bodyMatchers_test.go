package testHelper

import (
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

//  Native Body Matcher Tests
func TestNativeBodyMatcherNumber(t *testing.T) {
	expected := `4`
	var result int = 4
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherPrecision(t *testing.T) {
	expected := `4.11`
	var result float32 = 4.11
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherLong(t *testing.T) {
	expected := `411111111111111111`
	var result int64 = 411111111111111111
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherBoolean(t *testing.T) {
	expected := `true`
	var result bool = true
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherStringSlice(t *testing.T) {
	expected := `["Tuesday", "Saturday", "Wednesday", "Monday", "Sunday"]`
	var result []string = []string{
		"Tuesday", "Saturday", "Wednesday", "Monday", "Sunday",
	}
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherIntSlice(t *testing.T) {
	expected := `[1,2,3,4,5]`
	var result []int = []int{
		1, 2, 3, 4, 5,
	}
	NativeBodyMatcher(t, expected, result)
}

func TestNativeBodyMatcherBooleanError(t *testing.T) {
	expected := `nil`
	var result bool = true
	NativeBodyMatcher(&testing.T{}, expected, result)
}

//  Raw Body Matcher Tests
func TestRawBodyMatcherIntSlice(t *testing.T) {
	expected := `[1,2,3,4,5]`
	var result []int = []int{
		1, 2, 3, 4, 5,
	}
	RawBodyMatcher(t, expected, result)
}

func TestRawBodyMatcherBooleanError(t *testing.T) {
	expected := `nil`
	var result bool = true
	RawBodyMatcher(&testing.T{}, expected, result)
}

//  Is Same File Tests
func TestIsSameFile(t *testing.T) {
	expected := `https://www.gstatic.com/webp/gallery/1.jpg`
	result, err := https.GetFile(expected)
	if err != nil {
		t.Error("Error fetching File from ", expected)
	}
	IsSameAsFile(t, expected, result.File)
}

func TestIsSameFileError(t *testing.T) {
	expected := `http://localhost:3000/response/image`
	result, _ := https.GetFile("https://play.google.com/store/apps/dev?id=5700313618786177705&hl=en_US&gl=US")
	IsSameAsFile(&testing.T{}, expected, result.File)
}

func TestIsSameFileErrorURL(t *testing.T) {
	expected := `http://response/image`
	result, _ := https.GetFile(`http://localhost:3000/response/image`)
	IsSameAsFile(&testing.T{}, expected, result.File)
}

//  Slice to Comma Separated String Tests
func TestSliceToCommaSeparatedString(t *testing.T) {
	expected := `{"isMap": false,"id": "5a9fcb01caacc310dc6bab50"}`
	SliceToCommaSeparatedString(expected)
}

//  Keys And Values Body Matcher Tests
type Response struct {
	IsMap      bool       `json:"isMap"`
	Attributes Attributes `json:"attributes"`
	Id         string     `json:"id"`
}
type Attributes struct {
	Id string `json:"id"`
}

func TestKeysAndValuesBodyMatcherEmpty(t *testing.T) {
	expected := `{}`
	KeysAndValuesBodyMatcher[any](t, expected, nil, false, false)
}

func TestKeysAndValuesBodyMatcherObject(t *testing.T) {
	expected := `{"id": "5a9fcb01caacc310dc6bab51"}`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	KeysAndValuesBodyMatcher(t, expected, result, false, false)
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
	KeysAndValuesBodyMatcher(t, expected, result, false, false)
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
	KeysAndValuesBodyMatcher(&testing.T{}, expected, result, false, false)
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
	KeysAndValuesBodyMatcher(&testing.T{}, expected, result, false, false)
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
	KeysAndValuesBodyMatcher(&testing.T{}, expected, result, true, false)
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
	KeysAndValuesBodyMatcher(&testing.T{}, expected, result, true, false)
}

//  Keys Body Matcher Tests
func TestKeysBodyMatcherEmpty(t *testing.T) {
	expected := `{}`
	KeysBodyMatcher(t, expected, nil, false, false)
}

func TestKeysBodyMatcherObject(t *testing.T) {
	expected := `{"id": "5a9fcb01caacc310dc6bab51"}`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	KeysBodyMatcher(t, expected, result, false, false)
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
	KeysBodyMatcher(t, expected, result, false, false)
}

func TestKeysBodyMatcherObjectError(t *testing.T) {
	expected := `{"idd": "nil"}`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	KeysBodyMatcher(&testing.T{}, expected, result, false, false)
}

func TestKeysBodyMatcherUnmarshallingError(t *testing.T) {
	expected := `{"idd": "nil"`
	result := Attributes{
		Id: "5a9fcb01caacc310dc6bab51",
	}
	KeysBodyMatcher(&testing.T{}, expected, result, false, false)
}
