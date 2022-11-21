package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/apimatic/go-core-runtime/apiError"
)

const DEFAULT_DATE = "2006-01-02"

func PrepareQueryParams(queryParams url.Values, data map[string]interface{}) url.Values {
	if queryParams == nil {
		queryParams = url.Values{}
	}

	for k, v := range data {
		queryParams.Add(k, fmt.Sprintf("%v", v))
	}
	return queryParams
}

func JsonDecoderToString(dec *json.Decoder) (string, apiError.CustomError) {
	var customError apiError.CustomError
	var str string
	for {
		if err := dec.Decode(&str); err == io.EOF {
			break
		} else if err != nil {
			customError = *apiError.NewCustomError(log.Llongfile, "", err.Error(), err)
		}
	}
	return str, customError
}

func JsonDecoderToStringSlice(dec *json.Decoder) ([]string, apiError.CustomError) {
	var customError apiError.CustomError
	var arr []string
	for {
		if err := dec.Decode(&arr); err == io.EOF {
			break
		} else if err != nil {
			customError = *apiError.NewCustomError(log.Llongfile, "", err.Error(), err)
		}
	}
	return arr, customError
}

func JsonDecoderToIntSlice(dec *json.Decoder) ([]int, apiError.CustomError) {
	var customError apiError.CustomError
	var arr []int
	for {
		if err := dec.Decode(&arr); err == io.EOF {
			break
		} else if err != nil {
			customError = *apiError.NewCustomError(log.Llongfile, "", err.Error(), err)
		}
	}
	return arr, customError
}

func JsonDecoderToBooleanSlice(dec *json.Decoder) ([]bool, apiError.CustomError) {
	var customError apiError.CustomError
	var arr []bool
	for {
		if err := dec.Decode(&arr); err == io.EOF {
			break
		} else if err != nil {
			customError = *apiError.NewCustomError(log.Llongfile, "", err.Error(), err)
		}
	}
	return arr, customError
}

// ToTimeSlice is used to make a time.Time slice from a string slice.
func ToTimeSlice(slice interface{}, format string) ([]time.Time, apiError.CustomError) {
	var customError apiError.CustomError
	result := make([]time.Time, 0)
	if slice == nil {
		return []time.Time{}, customError
	}

	if format == time.UnixDate {
		for _, val := range slice.([]int64) {
			date := time.Unix(val, 0)
			result = append(result, date)
		}
	} else {
		for _, val := range slice.([]string) {
			date, err := time.Parse(format, val)
			if err != nil {
				customError = *apiError.NewCustomError(log.Llongfile, "", "Error parsing the date:", err)
			}
			result = append(result, date)
		}
	}
	return result, customError
}

// TimeToStringSlice is used to make a string slice from a time.Time slice.
func TimeToStringSlice(slice []time.Time, format string) []string {
	result := make([]string, 0)
	if slice == nil {
		return []string{}
	}

	for _, val := range slice {
		var date string
		if format == time.UnixDate {
			date = strconv.FormatInt(val.Unix(), 10)
		} else {
			date = val.Format(format)
		}
		result = append(result, date)
	}
	return result
}

// ToTimeMap is used to make a time.Time map from a string map.
func ToTimeMap(dict interface{}, format string) (map[string]time.Time, apiError.CustomError) {
	var customError apiError.CustomError
	result := make(map[string]time.Time, 0)
	if dict == nil {
		return map[string]time.Time{}, customError
	}

	if format == time.UnixDate {
		for key, val := range dict.(map[string]int64) {
			date := time.Unix(val, 0)
			result[key] = date
		}
	} else {
		for key, val := range dict.(map[string]string) {
			date, err := time.Parse(format, val)
			if err != nil {
				customError = *apiError.NewCustomError(log.Llongfile, "", "Error parsing the date: ", err)
			}
			result[key] = date
		}
	}
	return result, customError
}

// ToNullableTimeMap is used to make a nullable time.Time map from a string map.
func ToNullableTimeMap(dict interface{}, format string) (map[string]*time.Time, apiError.CustomError) {
	var customError apiError.CustomError
	result := make(map[string]*time.Time, 0)
	if dict == nil {
		return map[string]*time.Time{}, customError
	}

	if format == time.UnixDate {
		for key, val := range dict.(map[string]*int64) {
			if val == nil {
				result[key] = nil
			} else {
				date := time.Unix(*val, 0)
				result[key] = &date
			}
		}
	} else {
		for key, val := range dict.(map[string]*string) {
			if val == nil {
				result[key] = nil
			} else {
				date, err := time.Parse(format, *val)
				if err != nil {
					customError =*apiError.NewCustomError(log.Llongfile, "", "Error parsing the date: ", err)
				}
				result[key] = &date
			}
		}
	}
	return result, customError
}

// TimeToStringMap is used to make a string map from a time.Time map.
func TimeToStringMap(dict map[string]time.Time, format string) map[string]string {
	result := make(map[string]string, 0)
	if dict == nil {
		return map[string]string{}
	}

	for key, val := range dict {
		var date string
		if format == time.UnixDate {
			date = strconv.FormatInt(val.Unix(), 10)
		} else {
			date = val.Format(format)
		}
		result[key] = date
	}
	return result
}

// NullableTimeToStringMap is used to make a nullable string map from a time.Time map.
func NullableTimeToStringMap(dict map[string]*time.Time, format string) map[string]*string {
	result := make(map[string]*string, 0)
	if dict == nil {
		return map[string]*string{}
	}

	for key, val := range dict {
		if val == nil {
			result[key] = nil
		} else {
			var date string
			if format == time.UnixDate {
				date = strconv.FormatInt(val.Unix(), 10)
			} else {
				date = val.Format(format)
			}
			result[key] = &date
		}
	}
	return result
}

func UpdateUserAgent(userAgent string) string {
	updatedAgent := userAgent
	updatedAgent = strings.Replace(updatedAgent, "{os-info}", runtime.GOOS, -1)
	updatedAgent = strings.Replace(updatedAgent, "{engine}", runtime.Version(), -1)
	updatedAgent = strings.Replace(updatedAgent, "{engine-version}", strings.Replace(runtime.Version(), "go", "", 1), -1)

	return updatedAgent
}
