package https

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
)

// FormParam is a struct that represents a key-value pair for form parameters.
// It contains the key, value, and headers associated with the form parameter.
type FormParam struct {
	Key     string
	Value   any
	Headers http.Header
}

// structToMap converts a given data structure to a map.
func structToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	return mapData, err
}

// formEncodeMap recursively encodes the form parameters from a nested data structure.
// It processes the FormParam field and appends the result to the formParams slice.
// This function is used to handle nested data structures in form parameters.
func formEncodeMap(field FormParam, formParams *[]FormParam) ([]FormParam, error) {
	if formParams == nil {
		formParams = &[]FormParam{}
	}

	value := field.Value
	name := field.Key
	headers := field.Headers

	if value == nil {
		return *formParams, nil
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Ptr:
		structMap, err := structToMap(value)
		if err != nil {
			return nil, err
		}
		for key, value := range structMap {
			var fullName string = key
			if name != "" {
				fullName = name + "[" + key + "]"
			}
			formEncodeMap(FormParam{fullName, value, headers}, formParams)
		}
	case reflect.Map:
		for key, val := range value.(map[string]interface{}) {
			var fullName string = key
			if name != "" {
				fullName = name + "[" + key + "]"
			}
			formEncodeMap(FormParam{fullName, val, headers}, formParams)
		}
	case reflect.Slice:
		if reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			for num, val := range value.([]interface{}) {
				fullName := name + "[" + fmt.Sprintf("%v", num) + "]"
				formEncodeMap(FormParam{fullName, val, headers}, formParams)
			}
		} else {
			reflectValue := reflect.ValueOf(value)
			for num := 0; num < reflectValue.Len(); num++ {
				fullName := name + "[" + fmt.Sprintf("%v", num) + "]"
				formEncodeMap(FormParam{fullName, reflectValue.Index(num).Interface(), headers}, formParams)
			}
		}
	default:
		*formParams = append(*formParams, FormParam{name, value, headers})
	}
	return *formParams, nil
}

// prepareFormFields prepares the form fields from the given FormParam and adds them to the provided url.Values.
// It processes the FormParam field and encodes the value according to its data type.
func prepareFormFields(field FormParam, form url.Values) (url.Values, error) {
	if form == nil {
		form = url.Values{}
	}

	switch value := field.Value.(type) {
	case []string:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []int:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []int16:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []int32:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []int64:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []bool:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []float32:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	case []float64:
		for _, val := range value {
			form.Add(field.Key, fmt.Sprintf("%v", val))
		}
	default:
		formParams, err := formEncodeMap(field, nil)

		if err != nil {
			return nil, fmt.Errorf("Error parsing the date: %v", err)
		}
		for _, param := range formParams {
			form.Add(param.Key, fmt.Sprintf("%v", param.Value))
		}
	}

	return form, nil
}

// prepareMultipartFields prepares the multipart fields from the given FormParam and
// returns the body as a bytes.Buffer, along with the Content-Type header for the multipart form data.
func prepareMultipartFields(fields []FormParam) (bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, field := range fields {
		switch fieldValue := field.Value.(type) {
		case FileWrapper:
			headers := make(textproto.MIMEHeader)
			contentDisposition := mime.FormatMediaType("form-data", map[string]string{
				"name":     field.Key,
				"filename": fieldValue.FileName,
			})
			headers.Set("Content-Disposition", contentDisposition)
			if contentType := field.Headers.Get("Content-Type"); contentType != "" {
				headers.Set("Content-Type", contentType)
			}
			part, err := writer.CreatePart(headers)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			_, err = part.Write(fieldValue.File)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
		case string:
			headers := make(textproto.MIMEHeader)
			contentDisposition := mime.FormatMediaType("form-data", map[string]string{
				"name": field.Key,
			})
			headers.Set("Content-Disposition", contentDisposition)
			if contentType := field.Headers.Get("Content-Type"); contentType != "" {
				headers.Set("Content-Type", contentType)
			}
			part, err := writer.CreatePart(headers)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			_, err = part.Write([]byte(fieldValue))
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
		default:
			headers := make(textproto.MIMEHeader)
			contentDisposition := mime.FormatMediaType("form-data", map[string]string{
				"name": field.Key,
			})
			headers.Set("Content-Disposition", contentDisposition)
			if contentType := field.Headers.Get("Content-Type"); contentType != "" {
				headers.Set("Content-Type", contentType)
			}
			part, err := writer.CreatePart(headers)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			marshalledBytes, err := json.Marshal(fieldValue)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			_, err = part.Write(marshalledBytes)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
		}
	}
	writer.Close()
	return *body, writer.FormDataContentType(), nil
}
