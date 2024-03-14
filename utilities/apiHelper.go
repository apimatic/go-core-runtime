// Package utilities provides various utility functions and helpers for common operations.
// Copyright (c) APIMatic. All rights reserved.
package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_DATE = "2006-01-02"

// DecodeResults decodes JSON data from the provided json.Decoder into a given type T.
func DecodeResults[T any](decoder *json.Decoder) (T, error) {
	var result T
	for {
		if err := decoder.Decode(&result); err == io.EOF {
			break
		} else if err != nil {
			return result, err
		}
	}
	return result, nil
}

// PrepareQueryParams adds key-value pairs from the data map to the existing URL query parameters.
func PrepareQueryParams(queryParams url.Values, data map[string]any) url.Values {
	if queryParams == nil {
		queryParams = url.Values{}
	}

	for k, v := range data {
		queryParams.Add(k, fmt.Sprintf("%v", v))
	}
	return queryParams
}

// JsonDecoderToString decodes a JSON value from the provided json.Decoder into a string.
func JsonDecoderToString(dec *json.Decoder) (string, error) {
	return DecodeResults[string](dec)
}

// JsonDecoderToStringSlice decodes a JSON array from the provided json.Decoder into a string slice.
func JsonDecoderToStringSlice(dec *json.Decoder) ([]string, error) {
	return DecodeResults[[]string](dec)
}

// JsonDecoderToIntSlice decodes a JSON array from the provided json.Decoder into an int slice.
func JsonDecoderToIntSlice(dec *json.Decoder) ([]int, error) {
	return DecodeResults[[]int](dec)
}

// JsonDecoderToBooleanSlice decodes a JSON array from the provided json.Decoder into a bool slice.
func JsonDecoderToBooleanSlice(dec *json.Decoder) ([]bool, error) {
	return DecodeResults[[]bool](dec)
}

// ToTimeSlice converts a slice of strings or int64 values to a slice of time.Time values using the specified format.
func ToTimeSlice(slice any, format string) ([]time.Time, error) {
	result := make([]time.Time, 0)
	if slice == nil {
		return []time.Time{}, nil
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
				return nil, fmt.Errorf("error parsing the date: %v", err)
			}
			result = append(result, date)
		}
	}
	return result, nil
}

// TimeToStringSlice converts a slice of time.Time values to a slice of strings using the specified format.
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

// ToTimeMap converts a map with string or int64 values to a map with time.Time values using the specified format.
func ToTimeMap(dict any, format string) (map[string]time.Time, error) {
	result := make(map[string]time.Time)
	if dict == nil {
		return map[string]time.Time{}, nil
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
				return nil, fmt.Errorf("error parsing the date: %v", err)
			}
			result[key] = date
		}
	}
	return result, nil
}

// ToNullableTimeMap converts a map with nullable string or int64 values to a map with nullable time.Time values using the specified format.
func ToNullableTimeMap(dict any, format string) (map[string]*time.Time, error) {
	result := make(map[string]*time.Time)
	if dict == nil {
		return map[string]*time.Time{}, nil
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
					return nil, fmt.Errorf("error parsing the date: %v", err)
				}
				result[key] = &date
			}
		}
	}
	return result, nil
}

// TimeToStringMap converts a map with time.Time values to a map with strings using the specified format.
func TimeToStringMap(dict map[string]time.Time, format string) map[string]string {
	result := make(map[string]string)
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

// NullableTimeToStringMap converts a map with nullable time.Time values to a map with nullable strings using the specified format.
func NullableTimeToStringMap(dict map[string]*time.Time, format string) map[string]*string {
	result := make(map[string]*string)
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

// UpdateUserAgent replaces placeholders in the user agent string with the actual values.
func UpdateUserAgent(userAgent string) string {
	updatedAgent := userAgent
	updatedAgent = strings.Replace(updatedAgent, "{os-info}", runtime.GOOS, -1)
	updatedAgent = strings.Replace(updatedAgent, "{engine}", runtime.Version(), -1)
	updatedAgent = strings.Replace(updatedAgent, "{engine-version}", strings.Replace(runtime.Version(), "go", "", 1), -1)

	return updatedAgent
}
