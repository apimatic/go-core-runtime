package utilities

import (
	"bytes"
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

// NullableTimeToStringMap
func TestNullableTimeToStringMapNil(t *testing.T) {
	expected := map[string]*string{}
	result := NullableTimeToStringMap(nil, time.UnixDate)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestUnixNullableTimeToStringMap(t *testing.T) {
	input := GetNullableTimeMap(time.UnixDate)
	result := NullableTimeToStringMap(input, time.UnixDate)

	expected := make(map[string]*string)
	time1 := "1660992485"
	time2 := "1629456485"
	expected["time1"] = &time1
	expected["time2"] = &time2
	expected["time3"] = nil

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestRFC3339NullableTimeToStringMap(t *testing.T) {
	input := GetNullableTimeMap(time.RFC3339)
	result := NullableTimeToStringMap(input, time.RFC3339)

	expected := make(map[string]*string)
	time1 := "2022-08-20T15:48:05+05:00"
	time2 := "2021-08-20T15:48:05+05:00"
	expected["time1"] = &time1
	expected["time2"] = &time2
	expected["time3"] = nil

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestRFC1123NullableTimeToStringMap(t *testing.T) {
	input := GetNullableTimeMap(time.RFC1123)
	result := NullableTimeToStringMap(input, time.RFC1123)

	expected := make(map[string]*string)
	time1 := "Sat, 20 Aug 2022 15:48:05 PKT"
	time2 := "Fri, 20 Aug 2021 15:48:05 PKT"
	expected["time1"] = &time1
	expected["time2"] = &time2
	expected["time3"] = nil

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestDefaultNullableTimeToStringMap(t *testing.T) {
	input := GetNullableTimeMap(DEFAULT_DATE)
	result := NullableTimeToStringMap(input, DEFAULT_DATE)

	expected := make(map[string]*string)
	expectedtime1 := "2022-08-20"
	expectedtime2 := "2021-08-20"
	expected["time1"] = &expectedtime1
	expected["time2"] = &expectedtime2
	expected["time3"] = nil

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// TimeToStringMap
func TestTimeToStringMapNil(t *testing.T) {
	expected := map[string]string{}
	result := TimeToStringMap(nil, time.UnixDate)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestUnixTimeToStringMap(t *testing.T) {
	input := GetTimeMap(time.UnixDate)
	expected := map[string]string{"time1": "1660992485", "time2": "1629456485"}
	result := TimeToStringMap(input, time.UnixDate)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestRFC3339TimeToStringMap(t *testing.T) {
	input := GetTimeMap(time.RFC3339)
	expected := map[string]string{"time1": "2022-08-20T15:48:05+05:00", "time2": "2021-08-20T15:48:05+05:00"}
	result := TimeToStringMap(input, time.RFC3339)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestRFC1123TimeToStringMap(t *testing.T) {
	input := GetTimeMap(time.RFC1123)
	expected := map[string]string{"time1": "Sat, 20 Aug 2022 15:48:05 PKT", "time2": "Fri, 20 Aug 2021 15:48:05 PKT"}
	result := TimeToStringMap(input, time.RFC1123)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestDefaultTimeToStringMap(t *testing.T) {
	input := GetTimeMap(DEFAULT_DATE)
	expected := map[string]string{"time1": "2022-08-20", "time2": "2021-08-20"}
	result := TimeToStringMap(input, DEFAULT_DATE)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// ToNullableTimeMap
func TestToNullableTimeMapUnix(t *testing.T) {
	input := make(map[string]*int64)
	var time1 int64 = 1660992485
	var time2 int64 = 1629456485
	input["time1"] = &time1
	input["time2"] = &time2
	input["time3"] = nil

	result, _ := ToNullableTimeMap(input, time.UnixDate)
	expected := GetNullableTimeMap(time.UnixDate)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToNullableTimeMapRFC3339(t *testing.T) {
	input := make(map[string]*string)
	time1 := "2022-08-20T15:48:05+05:00"
	time2 := "2021-08-20T15:48:05+05:00"
	input["time1"] = &time1
	input["time2"] = &time2
	input["time3"] = nil

	result, _ := ToNullableTimeMap(input, time.RFC3339)
	expected := GetNullableTimeMap(time.RFC3339)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToNullableTimeMapParsingError(t *testing.T) {
	input := make(map[string]*string)
	time1 := "2022-08-20T15:48:05+05:00"
	time2 := "2021-08-20T"
	input["time1"] = &time1
	input["time2"] = &time2
	input["time3"] = nil

	result, err := ToNullableTimeMap(input, time.RFC3339)
	if err == nil {
		t.Errorf("The code should get error while parsing date time.")
		expected := GetNullableTimeMap(time.RFC3339)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestToNullableTimeMapRFC1123(t *testing.T) {
	input := make(map[string]*string)
	time1 := "Sat, 20 Aug 2022 15:48:05 PKT"
	time2 := "Fri, 20 Aug 2021 15:48:05 PKT"
	input["time1"] = &time1
	input["time2"] = &time2
	input["time3"] = nil

	result, _ := ToNullableTimeMap(input, time.RFC1123)

	expected := GetNullableTimeMap(time.RFC1123)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToNullableTimeMapDefault(t *testing.T) {
	input := make(map[string]*string)
	time1 := "2022-08-20"
	time2 := "2021-08-20"
	input["time1"] = &time1
	input["time2"] = &time2
	input["time3"] = nil

	result, _ := ToNullableTimeMap(input, DEFAULT_DATE)

	expected := GetNullableTimeMap(DEFAULT_DATE)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToNullableTimeMapNil(t *testing.T) {
	result, _ := ToNullableTimeMap(nil, DEFAULT_DATE)
	expected := map[string]*time.Time{}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// ToTimeMap
func TestToTimeMapNil(t *testing.T) {
	result, _ := ToTimeMap(nil, time.UnixDate)

	expected := map[string]time.Time{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeMapParsingError(t *testing.T) {
	input := map[string]string{
		"time1": "2022-08-20T15:48:05+05:00",
		"time2": "2021-08-20T15:48:05",
	}
	result, err := ToTimeMap(input, time.RFC3339)
	if err == nil {
		t.Errorf("The code should get error while parsing date time.")
		expected := GetTimeMap(time.RFC3339)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestToTimeMapUnix(t *testing.T) {
	input := map[string]int64{
		"time1": 1660992485,
		"time2": 1629456485,
	}
	result, _ := ToTimeMap(input, time.UnixDate)

	expected := GetTimeMap(time.UnixDate)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeMapRFC3339(t *testing.T) {
	input := map[string]string{
		"time1": "2022-08-20T15:48:05+05:00",
		"time2": "2021-08-20T15:48:05+05:00",
	}
	result, _ := ToTimeMap(input, time.RFC3339)

	expected := GetTimeMap(time.RFC3339)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeMapRFC1123(t *testing.T) {
	input := map[string]string{
		"time1": "Sat, 20 Aug 2022 15:48:05 PKT",
		"time2": "Fri, 20 Aug 2021 15:48:05 PKT",
	}
	result, _ := ToTimeMap(input, time.RFC1123)

	expected := GetTimeMap(time.RFC1123)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeMapDefault(t *testing.T) {
	input := map[string]string{
		"time1": "2022-08-20",
		"time2": "2021-08-20",
	}
	result, _ := ToTimeMap(input, DEFAULT_DATE)

	expected := GetTimeMap(DEFAULT_DATE)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// TimeToStringSlice
func TestTimeToStringSliceNil(t *testing.T) {
	expected := []string{}
	result := TimeToStringSlice(nil, time.UnixDate)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestUnixTimeToStringSlice(t *testing.T) {
	slice := GetTimeSlice(time.UnixDate)
	expected := []string{"1660992485", "1629456485"}
	result := TimeToStringSlice(slice, time.UnixDate)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestRFC3339TimeToStringSlice(t *testing.T) {
	slice := GetTimeSlice(time.RFC3339)
	expected := []string{"2022-08-20T15:48:05+05:00", "2021-08-20T15:48:05+05:00"}
	result := TimeToStringSlice(slice, time.RFC3339)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestRFC1123TimeToStringSlice(t *testing.T) {
	slice := GetTimeSlice(time.RFC1123)
	expected := []string{"Sat, 20 Aug 2022 15:48:05 PKT", "Fri, 20 Aug 2021 15:48:05 PKT"}
	result := TimeToStringSlice(slice, time.RFC1123)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestDefaultTimeToStringSlice(t *testing.T) {
	slice := GetTimeSlice(DEFAULT_DATE)
	expected := []string{"2022-08-20", "2021-08-20"}
	result := TimeToStringSlice(slice, DEFAULT_DATE)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// ToTimeSlice
func TestToTimeSliceNil(t *testing.T) {
	result, _ := ToTimeSlice(nil, time.UnixDate)

	expected := []time.Time{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeSliceParsingError(t *testing.T) {
	slice := []string{"2022-08-20T15:48:05+05:00", "2021-08-20T"}
	result, err := ToTimeSlice(slice, time.RFC3339)

	if err == nil {
		t.Errorf("The code should get error while parsing date time.")
		expected := GetTimeSlice(time.RFC3339)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestToTimeSliceUnix(t *testing.T) {
	slice := []int64{1660992485, 1629456485}
	result, _ := ToTimeSlice(slice, time.UnixDate)

	expected := GetTimeSlice(time.UnixDate)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeSliceRFC3339(t *testing.T) {
	slice := []string{"2022-08-20T15:48:05+05:00", "2021-08-20T15:48:05+05:00"}
	result, _ := ToTimeSlice(slice, time.RFC3339)

	expected := GetTimeSlice(time.RFC3339)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeSliceRFC1123(t *testing.T) {
	slice := []string{"Sat, 20 Aug 2022 15:48:05 PKT", "Fri, 20 Aug 2021 15:48:05 PKT"}
	result, _ := ToTimeSlice(slice, time.RFC1123)

	expected := GetTimeSlice(time.RFC1123)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestToTimeSliceDefault(t *testing.T) {
	slice := []string{"2022-08-20", "2021-08-20"}
	result, _ := ToTimeSlice(slice, DEFAULT_DATE)

	time1, _ := time.Parse(DEFAULT_DATE, "2022-08-20")
	time2, _ := time.Parse(DEFAULT_DATE, "2021-08-20")
	expected := []time.Time{time1, time2}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// JsonDecoderToBooleanSlice
func TestJsonDecoderToBooleanSliceError(t *testing.T) {
	boolSlice := []bool{true, false}
	result, err := JsonDecoderToBooleanSlice(GetJsonDecoded([]int{1, 2}))

	if err == nil {
		t.Errorf("The code should get error while decoding.")
		expected := boolSlice
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestJsonDecoderToBooleanSlice(t *testing.T) {
	boolSlice := []bool{true, false}
	result, _ := JsonDecoderToBooleanSlice(GetJsonDecoded(boolSlice))

	expected := boolSlice
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestJsonDecoderToBooleanSliceWithEmptySlice(t *testing.T) {
	boolSlice := []bool{}
	result, _ := JsonDecoderToBooleanSlice(GetJsonDecoded(boolSlice))

	expected := boolSlice
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// JsonDecoderToIntSlice
func TestJsonDecoderToIntSliceError(t *testing.T) {
	intSlice := []int{1, 2}
	result, err := JsonDecoderToIntSlice(GetJsonDecoded([]bool{true, false}))
	if err == nil {
		t.Errorf("The code should get error while decoding.")
		expected := intSlice
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestJsonDecoderToIntSlice(t *testing.T) {
	intSlice := []int{1, 2}
	result, _ := JsonDecoderToIntSlice(GetJsonDecoded(intSlice))

	expected := intSlice
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestJsonDecoderToIntSliceWithEmptySlice(t *testing.T) {
	intSlice := []int{}
	result, _ := JsonDecoderToIntSlice(GetJsonDecoded(intSlice))

	expected := intSlice
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// JsonDecoderToStringSlice
func TestJsonDecoderToStringSliceError(t *testing.T) {
	stringSlice := []string{"GO", "APIMatic"}
	result, err := JsonDecoderToStringSlice(GetJsonDecoded([]bool{true, false}))

	if err == nil {
		t.Errorf("The code should get error while decoding.")
		expected := stringSlice
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestJsonDecoderToStringSlice(t *testing.T) {
	stringSlice := []string{"GO", "APIMatic"}
	result, _ := JsonDecoderToStringSlice(GetJsonDecoded(stringSlice))

	expected := stringSlice
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestJsonDecoderToStringSliceWithEmptySlice(t *testing.T) {
	stringSlice := []string{}
	result, _ := JsonDecoderToStringSlice(GetJsonDecoded(stringSlice))

	expected := stringSlice
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// JsonDecoderToString
func TestJsonDecoderToStringError(t *testing.T) {
	result, err := JsonDecoderToString(GetJsonDecoded(34))

	if err == nil {
		t.Errorf("The code should get error while decoding.")
		expected := "This is Core Library for Go."
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
		}
	}
}

func TestJsonDecoderToString(t *testing.T) {
	result, _ := JsonDecoderToString(GetJsonDecoded("This is Core Library for Go."))

	expected := "This is Core Library for Go."
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestJsonDecoderToStringWithEmptyString(t *testing.T) {
	result, _ := JsonDecoderToString(GetJsonDecoded(""))

	expected := ""
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// PrepareQueryParams
func TestPrepareQueryParamsDuplicateData(t *testing.T) {
	result := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	data := map[string]interface{}{
		"key":  "value",
		"key1": 1,
	}
	PrepareQueryParams(result, data)
	expected := url.Values{
		"key":  []string{"value", "value"},
		"key1": []string{"1", "1"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareQueryParamsNilData(t *testing.T) {
	result := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}

	PrepareQueryParams(result, nil)
	expected := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareQueryParamsNilQueryParams(t *testing.T) {
	data := map[string]interface{}{
		"key":  "value",
		"key1": 1,
	}
	result := url.Values{}
	PrepareQueryParams(result, data)
	expected := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareQueryParamsEmptyQueryParams(t *testing.T) {
	result := url.Values{}
	data := map[string]interface{}{
		"key":  "value",
		"key1": 1,
	}
	PrepareQueryParams(result, data)
	expected := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareQueryParamsAppendQueryParams(t *testing.T) {
	result := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	data := map[string]interface{}{
		"key":  "value1",
		"key1": 2,
	}
	PrepareQueryParams(result, data)
	expected := url.Values{
		"key":  []string{"value", "value1"},
		"key1": []string{"1", "2"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestPrepareQueryParamsAppendEmptyData(t *testing.T) {
	result := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	data := map[string]interface{}{}
	PrepareQueryParams(result, data)
	expected := url.Values{
		"key":  []string{"value"},
		"key1": []string{"1"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// UpdateUserAgent
func TestUpdateUserAgentAllArguments(t *testing.T) {
	result := UpdateUserAgent("userAgent {os-info} {engine} {engine-version}")
	if !strings.Contains(result, "userAgent linux go") {
		t.Error("Incorrect UserAgent. Got:", result)
	}
}

func TestUpdateUserAgentEmptyArguments(t *testing.T) {
	result := UpdateUserAgent("userAgent")
	if result != "userAgent" {
		t.Error("Fails")
	}
}

func TestUpdateUserAgent2Arguments(t *testing.T) {
	result := UpdateUserAgent("userAgent {os-info} {engine}")
	if !strings.Contains(result, "userAgent linux go") {
		t.Error("Incorrect UserAgent. Got:", result)
	}
}

func TestUpdateUserAgentWrongArguments(t *testing.T) {
	result := UpdateUserAgent("userAgent {info} {enginee}")
	if result != "userAgent {info} {enginee}" {
		t.Error("Fails")
	}
}

func TestDecodeResultsString(t *testing.T) {
	expected := "This is Core Library for Go."
	decoder := GetJsonDecoded(expected)
	result, _ := DecodeResults[string](decoder)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

func TestDecodeResultsInt(t *testing.T) {
	expected := "This is Core Library for Go."
	decoder := GetJsonDecoded(expected)
	result, _ := DecodeResults[int](decoder)
	if reflect.DeepEqual(result, expected) {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result)
	}
}

// Helper methods
func GetTimeMap(format string) map[string]time.Time {
	var time1, time2 time.Time
	if format == time.RFC3339 {
		time1, _ = time.Parse(time.RFC3339, "2022-08-20T15:48:05+05:00")
		time2, _ = time.Parse(time.RFC3339, "2021-08-20T15:48:05+05:00")
	} else if format == time.RFC1123 {
		time1, _ = time.Parse(time.RFC1123, "Sat, 20 Aug 2022 15:48:05 PKT")
		time2, _ = time.Parse(time.RFC1123, "Fri, 20 Aug 2021 15:48:05 PKT")
	} else if format == time.UnixDate {
		time1 = time.Unix(1660992485, 0)
		time2 = time.Unix(1629456485, 0)
	} else if format == DEFAULT_DATE {
		time1, _ = time.Parse(DEFAULT_DATE, "2022-08-20")
		time2, _ = time.Parse(DEFAULT_DATE, "2021-08-20")
	}

	return map[string]time.Time{"time1": time1, "time2": time2}
}

func GetNullableTimeMap(format string) map[string]*time.Time {
	var time1, time2 time.Time
	if format == time.RFC3339 {
		time1, _ = time.Parse(time.RFC3339, "2022-08-20T15:48:05+05:00")
		time2, _ = time.Parse(time.RFC3339, "2021-08-20T15:48:05+05:00")
	} else if format == time.RFC1123 {
		time1, _ = time.Parse(time.RFC1123, "Sat, 20 Aug 2022 15:48:05 PKT")
		time2, _ = time.Parse(time.RFC1123, "Fri, 20 Aug 2021 15:48:05 PKT")
	} else if format == time.UnixDate {
		time1 = time.Unix(1660992485, 0)
		time2 = time.Unix(1629456485, 0)
	} else if format == DEFAULT_DATE {
		time1, _ = time.Parse(DEFAULT_DATE, "2022-08-20")
		time2, _ = time.Parse(DEFAULT_DATE, "2021-08-20")
	}

	nullableMap := make(map[string]*time.Time)
	nullableMap["time1"] = &time1
	nullableMap["time2"] = &time2
	nullableMap["time3"] = nil

	return nullableMap
}

func GetTimeSlice(format string) []time.Time {
	var time1, time2 time.Time
	if format == time.RFC3339 {
		time1, _ = time.Parse(time.RFC3339, "2022-08-20T15:48:05+05:00")
		time2, _ = time.Parse(time.RFC3339, "2021-08-20T15:48:05+05:00")
	} else if format == time.RFC1123 {
		time1, _ = time.Parse(time.RFC1123, "Sat, 20 Aug 2022 15:48:05 PKT")
		time2, _ = time.Parse(time.RFC1123, "Fri, 20 Aug 2021 15:48:05 PKT")
	} else if format == time.UnixDate {
		time1 = time.Unix(1660992485, 0)
		time2 = time.Unix(1629456485, 0)
	} else if format == DEFAULT_DATE {
		time1, _ = time.Parse(DEFAULT_DATE, "2022-08-20")
		time2, _ = time.Parse(DEFAULT_DATE, "2021-08-20")
	}

	return []time.Time{time1, time2}
}

func GetJsonDecoded(arr interface{}) *json.Decoder {
	buffer := &bytes.Buffer{}
	json.NewEncoder(buffer).Encode(arr)
	byteSlice := buffer.Bytes()

	return json.NewDecoder(bytes.NewReader(byteSlice))
}
