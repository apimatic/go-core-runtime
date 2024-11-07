package utilities_test

import (
	"encoding/json"
	"github.com/apimatic/go-core-runtime/internal/assert"
	"github.com/apimatic/go-core-runtime/utilities"
	"testing"
	"time"
)

func TestUnixDateTimeString(t *testing.T) {
	newTime := time.Unix(1484719381, 0)
	unixDateTime := utilities.NewUnixDateTime(newTime)
	assert.Equal(t, "1484719381", unixDateTime.String())
}

func TestUnixDateTime(t *testing.T) {
	newTime := time.Unix(1484719381, 0)
	unixDateTime := utilities.NewUnixDateTime(newTime)

	bytes, err := json.Marshal(unixDateTime)
	assert.NoError(t, err)
	var newUnixDateTime utilities.UnixDateTime
	err = json.Unmarshal(bytes, &newUnixDateTime)
	assert.NoError(t, err)
	assert.Equal(t, newTime, newUnixDateTime.Value())
}

func TestUnixDateTimeError(t *testing.T) {
	var newUnixDateTime utilities.UnixDateTime
	err := json.Unmarshal([]byte(`"Sun, 06 Nov 1994 08:49:37 GMT"`), &newUnixDateTime)
	assert.EqualError(t, err, "json: cannot unmarshal string into Go value of type int64")
}

func TestDefaultTimeString(t *testing.T) {
	newTime, err := time.Parse(utilities.DEFAULT_DATE, "1994-02-13")
	assert.NoError(t, err)
	defaultTime := utilities.NewDefaultTime(newTime)
	assert.Equal(t, "1994-02-13", defaultTime.String())
}

func TestDefaultTime(t *testing.T) {
	newTime, err := time.Parse(utilities.DEFAULT_DATE, "1994-02-13")
	assert.NoError(t, err)
	defaultTime := utilities.NewDefaultTime(newTime)

	bytes, err := json.Marshal(defaultTime)
	assert.NoError(t, err)
	var newDefaultTime utilities.DefaultTime
	err = json.Unmarshal(bytes, &newDefaultTime)
	assert.NoError(t, err)
	assert.Equal(t, newTime, newDefaultTime.Value())
}

func TestDefaultTimeError1(t *testing.T) {
	var newDefaultTime utilities.DefaultTime
	err := json.Unmarshal([]byte(`1484719381`), &newDefaultTime)
	assert.EqualError(t, err, "json: cannot unmarshal number into Go value of type string")
}

func TestDefaultTimeError2(t *testing.T) {
	var newDefaultTime utilities.DefaultTime
	err := json.Unmarshal([]byte(`"Sun, 06 Nov 1994 08:49:37 GMT"`), &newDefaultTime)
	assert.EqualError(t, err, "parsing time \"Sun, 06 Nov 1994 08:49:37 GMT\" as \"2006-01-02\": cannot parse \"Sun, 06 Nov 1994 08:49:37 GMT\" as \"2006\"")
}

func TestRFC3339TimeString(t *testing.T) {
	newTime, err := time.Parse(time.RFC3339Nano, "1994-02-13T14:01:54.9571247Z")
	assert.NoError(t, err)
	rFC3339Time := utilities.NewRFC3339Time(newTime)
	assert.Equal(t, "1994-02-13T14:01:54.9571247Z", rFC3339Time.String())
}

func TestRFC3339Time(t *testing.T) {
	newTime, err := time.Parse(time.RFC3339Nano, "1994-02-13T14:01:54.9571247Z")
	assert.NoError(t, err)
	rFC3339Time := utilities.NewRFC3339Time(newTime)

	bytes, err := json.Marshal(rFC3339Time)
	assert.NoError(t, err)
	var newRFC3339Time utilities.RFC3339Time
	err = json.Unmarshal(bytes, &newRFC3339Time)
	assert.NoError(t, err)
	assert.Equal(t, newTime, newRFC3339Time.Value())
}

func TestRFC1123TimeString(t *testing.T) {
	newTime, err := time.Parse(time.RFC1123, "Sun, 06 Nov 1994 08:49:37 GMT")
	assert.NoError(t, err)
	rFC1123Time := utilities.NewRFC1123Time(newTime)
	assert.Equal(t, "Sun, 06 Nov 1994 08:49:37 GMT", rFC1123Time.String())
}

func TestRFC1123Time(t *testing.T) {
	newTime, err := time.Parse(time.RFC1123, "Sun, 06 Nov 1994 08:49:37 GMT")
	assert.NoError(t, err)
	rFC1123Time := utilities.NewRFC1123Time(newTime)

	bytes, err := json.Marshal(rFC1123Time)
	assert.NoError(t, err)
	var newRFC1123Time utilities.RFC1123Time
	err = json.Unmarshal(bytes, &newRFC1123Time)
	assert.NoError(t, err)
	assert.Equal(t, newTime, rFC1123Time.Value())
}

func TestObjSliceToTimeSlice(t *testing.T) {
	initialBytes := []byte(`["Sun, 06 Nov 1994 08:49:37 GMT","Sun, 06 Nov 1994 08:49:37 GMT"]`)
	var strArray []string
	err := json.Unmarshal(initialBytes, &strArray)
	assert.NoError(t, err)
	timeArray, err := utilities.ToTimeSlice(strArray, time.RFC1123)
	assert.NoError(t, err)
	newObjArray := utilities.TimeSliceToObjSlice[utilities.RFC1123Time](timeArray)
	newTimeArray := utilities.ObjSliceToTimeSlice(newObjArray)
	newStrArray := utilities.TimeToStringSlice(newTimeArray, time.RFC1123)
	resultBytes, err := json.Marshal(newStrArray)
	assert.NoError(t, err)
	assert.Equal(t, string(initialBytes), string(resultBytes))
}

func TestObjMapToTimeMap(t *testing.T) {
	initialBytes := []byte(`{"key1":"Sun, 06 Nov 1994 08:49:37 GMT","key2":"Sun, 06 Nov 1994 08:49:37 GMT"}`)
	var strMap map[string]string
	err := json.Unmarshal(initialBytes, &strMap)
	assert.NoError(t, err)
	timeMap, err := utilities.ToTimeMap(strMap, time.RFC1123)
	assert.NoError(t, err)
	newObjMap := utilities.TimeMapToObjMap[utilities.RFC1123Time](timeMap)
	newTimeMap := utilities.ObjMapToTimeMap(newObjMap)
	newStrMap := utilities.TimeToStringMap(newTimeMap, time.RFC1123)
	resultBytes, err := json.Marshal(newStrMap)
	assert.NoError(t, err)
	assert.Equal(t, string(initialBytes), string(resultBytes))
}
