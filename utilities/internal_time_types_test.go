package utilities

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnixDateTimeString(t *testing.T) {
	newTime := time.Unix(1484719381, 0)
	unixDateTime := NewUnixDateTime(newTime)
	AssertEquals(t, "1484719381", unixDateTime.String())
}

func TestUnixDateTime(t *testing.T) {
	newTime := time.Unix(1484719381, 0)
	unixDateTime := NewUnixDateTime(newTime)

	bytes, err := json.Marshal(unixDateTime)
	AssertNoError(t, err)
	var newUnixDateTime UnixDateTime
	err = json.Unmarshal(bytes, &newUnixDateTime)
	AssertNoError(t, err)
	AssertEquals(t, newTime, newUnixDateTime.Value())
}

func TestUnixDateTimeError(t *testing.T) {
	var newUnixDateTime UnixDateTime
	json.Unmarshal([]byte(`"Sun, 06 Nov 1994 08:49:37 GMT"`), &newUnixDateTime)
}

func TestDefaultTimeString(t *testing.T) {
	newTime, err := time.Parse(DEFAULT_DATE, "1994-02-13")
	AssertNoError(t, err)
	defaultTime := NewDefaultTime(newTime)
	AssertEquals(t, "1994-02-13", defaultTime.String())
}

func TestDefaultTime(t *testing.T) {
	newTime, err := time.Parse(DEFAULT_DATE, "1994-02-13")
	AssertNoError(t, err)
	defaultTime := NewDefaultTime(newTime)

	bytes, err := json.Marshal(defaultTime)
	AssertNoError(t, err)
	var newDefaultTime DefaultTime
	err = json.Unmarshal(bytes, &newDefaultTime)
	AssertNoError(t, err)
	AssertEquals(t, newTime, newDefaultTime.Value())
}

func TestDefaultTimeError1(t *testing.T) {
	var newDefaultTime DefaultTime
	json.Unmarshal([]byte(`1484719381`), &newDefaultTime)
}

func TestDefaultTimeError2(t *testing.T) {
	var newDefaultTime DefaultTime
	json.Unmarshal([]byte(`"Sun, 06 Nov 1994 08:49:37 GMT"`), &newDefaultTime)
}

func TestRFC3339TimeString(t *testing.T) {
	newTime, err := time.Parse(time.RFC3339Nano, "1994-02-13T14:01:54.9571247Z")
	AssertNoError(t, err)
	rFC3339Time := NewRFC3339Time(newTime)
	AssertEquals(t, "1994-02-13T14:01:54.9571247Z", rFC3339Time.String())
}

func TestRFC3339Time(t *testing.T) {
	newTime, err := time.Parse(time.RFC3339Nano, "1994-02-13T14:01:54.9571247Z")
	AssertNoError(t, err)
	rFC3339Time := NewRFC3339Time(newTime)

	bytes, err := json.Marshal(rFC3339Time)
	AssertNoError(t, err)
	var newRFC3339Time RFC3339Time
	err = json.Unmarshal(bytes, &newRFC3339Time)
	AssertNoError(t, err)
	AssertEquals(t, newTime, newRFC3339Time.Value())
}

func TestRFC1123TimeString(t *testing.T) {
	newTime, err := time.Parse(time.RFC1123, "Sun, 06 Nov 1994 08:49:37 GMT")
	AssertNoError(t, err)
	rFC1123Time := NewRFC1123Time(newTime)
	AssertEquals(t, "Sun, 06 Nov 1994 08:49:37 GMT", rFC1123Time.String())
}

func TestRFC1123Time(t *testing.T) {
	newTime, err := time.Parse(time.RFC1123, "Sun, 06 Nov 1994 08:49:37 GMT")
	AssertNoError(t, err)
	rFC1123Time := NewRFC1123Time(newTime)

	bytes, err := json.Marshal(rFC1123Time)
	AssertNoError(t, err)
	var newRFC1123Time RFC1123Time
	err = json.Unmarshal(bytes, &newRFC1123Time)
	AssertNoError(t, err)
	AssertEquals(t, newTime, rFC1123Time.Value())
}

func TestObjSliceToTimeSlice(t *testing.T) {
	initialBytes := []byte(`["Sun, 06 Nov 1994 08:49:37 GMT","Sun, 06 Nov 1994 08:49:37 GMT"]`)
	var strArray []string
	err := json.Unmarshal(initialBytes, &strArray)
	AssertNoError(t, err)
	timeArray, err := ToTimeSlice(strArray, time.RFC1123)
	AssertNoError(t, err)
	newObjArray := TimeSliceToObjSlice[RFC1123Time](timeArray)
	newTimeArray := ObjSliceToTimeSlice(newObjArray)
	newStrArray := TimeToStringSlice(newTimeArray, time.RFC1123)
	resultBytes, err := json.Marshal(newStrArray)
	AssertNoError(t, err)
	AssertEquals(t, string(initialBytes), string(resultBytes))
}

func TestObjMapToTimeMap(t *testing.T) {
	initialBytes := []byte(`{"key1":"Sun, 06 Nov 1994 08:49:37 GMT","key2":"Sun, 06 Nov 1994 08:49:37 GMT"}`)
	var strMap map[string]string
	err := json.Unmarshal(initialBytes, &strMap)
	AssertNoError(t, err)
	timeMap, err := ToTimeMap(strMap, time.RFC1123)
	AssertNoError(t, err)
	newObjMap := TimeMapToObjMap[RFC1123Time](timeMap)
	newTimeMap := ObjMapToTimeMap(newObjMap)
	newStrMap := TimeToStringMap(newTimeMap, time.RFC1123)
	resultBytes, err := json.Marshal(newStrMap)
	AssertNoError(t, err)
	AssertEquals(t, string(initialBytes), string(resultBytes))
}
