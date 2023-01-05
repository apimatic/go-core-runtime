package https

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strings"
)

func structToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	return mapData, err
}

func formEncodeMap(name string, value interface{}, keys *[]map[string]interface{}) ([]map[string]interface{}, error) {
	if keys == nil {
		keys = &[]map[string]interface{}{}
	}

	if value == nil {
		return append(*keys, make(map[string]interface{})), nil
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Ptr:
		structMap, err := structToMap(value)
		if err != nil {
			return nil, err
		}
		for k, v := range structMap {
			var fullName string = k
			if name != "" {
				fullName = name + "[" + k + "]"
			}
			formEncodeMap(fullName, v, keys)
		}
	case reflect.Map:
		for k, v := range value.(map[string]interface{}) {
			var fullName string = k
			if name != "" {
				fullName = name + "[" + k + "]"
			}
			formEncodeMap(fullName, v, keys)
		}
	case reflect.Slice:
		if reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
			for num, val := range value.([]interface{}) {
				fullName := name + "[" + fmt.Sprintf("%v", num) + "]"
				formEncodeMap(fullName, val, keys)
			}
		} else {
			reflectValue := reflect.ValueOf(value)
			for num := 0; num < reflectValue.Len(); num++ {
				fullName := name + "[" + fmt.Sprintf("%v", num) + "]"
				formEncodeMap(fullName, reflectValue.Index(num).Interface(), keys)
			}
		}
	default:
		*keys = append(*keys, map[string]interface{}{name: fmt.Sprintf("%v", value)})
	}

	return *keys, nil
}

func PrepareFormFields(key string, value interface{}, form url.Values) (url.Values, error) {
	if form == nil {
		form = url.Values{}
	}

	switch x := value.(type) {
	case []string:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []int:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []int16:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []int32:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []int64:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []bool:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []float32:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	case []float64:
		for _, val := range x {
			form.Add(key, fmt.Sprintf("%v", val))
		}
	default:
		k, err := formEncodeMap(key, value, nil)

		if err != nil {
			return nil, fmt.Errorf("Error parsing the date: %v", err)
		}
		for num := range k {
			for key, val := range k[num] {
				form.Add(key, fmt.Sprintf("%v", val))
			}
		}
	}

	return form, nil
}

func PrepareMultipartFields(fields map[string]interface{}) (bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range fields {
		switch x := val.(type) {
		case FileWrapper:
			fw, err := writer.CreateFormFile(key, x.FileName)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			_, err = io.Copy(fw, bytes.NewReader(x.File))
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
		default:
			fw, err := writer.CreateFormField(key)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			marshalledBytes, err := json.Marshal(x)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			_, err = io.Copy(fw, strings.NewReader(string(marshalledBytes)))
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
		}
	}
	writer.Close()
	return *body, writer.FormDataContentType(), nil
}
