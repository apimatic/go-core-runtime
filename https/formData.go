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

// toMap converts a FormParam to a map of string key-value pairs.
func (param *FormParam) toMap() (map[string][]string, error) {

	paramMap := make(map[string][]string)
	if param.Value == nil {
		return paramMap, nil
	}
	valueType := reflect.TypeOf(param.Value).Kind()
	switch valueType {
	case reflect.Map, reflect.Struct, reflect.Ptr:
		// Convert Struct and Pointer types into Map.
		if valueType == reflect.Struct || valueType == reflect.Ptr {
			structMap, err := structToMap(param.Value)
			if err != nil {
				return paramMap, err
			}
			param.Value = structMap
		}
		// Add Map key-value pairs into the parent Map.
		iter := reflect.ValueOf(param.Value).MapRange()
		for iter.Next() {
			innerParam := &FormParam{
				fmt.Sprintf("%v[%v]", param.Key, iter.Key()),
				iter.Value().Interface(),
				param.Headers,
			}
			innerParamMap, err := innerParam.toMap()
			if err != nil {
				return paramMap, err
			}
			for k, v := range innerParamMap {
				paramMap[k] = v
			}
		}
	case reflect.Slice:
		reflectValue := reflect.ValueOf(param.Value)
		for index := 0; index < reflectValue.Len(); index++ {
			//frmt := "%v[%v]"
			innerParam := &FormParam{
				fmt.Sprintf("%v[%v]", param.Key, index),
				reflectValue.Index(index).Interface(),
				param.Headers,
			}
			innerParamMap, err := innerParam.toMap()
			if err != nil {
				return paramMap, err
			}
			for _, values := range innerParamMap {
				for _, value := range values{
					paramMap[param.Key] = append(paramMap[param.Key], value)
				}
			}
		}
	case reflect.String:
		paramMap[param.Key] = append(paramMap[param.Key], param.Value.(string))
	case reflect.Array:
		paramMap[param.Key] = append(paramMap[param.Key], fmt.Sprintf("%v", param.Value))
	default:
		bytes, err := json.Marshal(param.Value)
		if err == nil {
			paramMap[param.Key] =  append(paramMap[param.Key], string(bytes))
		} else {
			paramMap[param.Key] = append(paramMap[param.Key], fmt.Sprintf("%v", param.Value))
		}
	}
	return paramMap, nil
}

// FormParams represents a collection of FormParam objects.
type FormParams []FormParam

// Add appends a FormParam to the FormParams collection.
func (fp *FormParams) Add(formParam FormParam) {
	if formParam.Value != nil {
		*fp = append(*fp, formParam)
	}
}

// prepareFormFields prepares the form fields from the given FormParams and adds them to the url.Values.
// It processes each FormParam field and encodes the value according to its data type.
func (fp *FormParams) prepareFormFields(form url.Values) error {
	if form == nil {
		form = url.Values{}
	}
	for _, param := range *fp {
		paramsMap, err := param.toMap()
		if err != nil {
			return err
		}
		for key, values := range paramsMap {
			for _, value := range values {
				form.Add(key, value)
			}
		}
	}
	return nil
}

// prepareMultipartFields prepares the multipart fields from the given FormParams and
// returns the body as a bytes.Buffer, along with the Content-Type header for the multipart form data.
func (fp *FormParams) prepareMultipartFields() (bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, field := range *fp {
		switch fieldValue := field.Value.(type) {
		case FileWrapper:
			mediaParam := map[string]string{
				"name":     field.Key,
				"filename": fieldValue.FileName,
			}
			formParamWriter(writer, field.Headers, mediaParam, fieldValue.File)
		default:
			paramsMap, err := field.toMap()
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			for key, values := range paramsMap {
				mediaParam := map[string]string{"name": key}
				for _, value := range values {
					formParamWriter(writer, field.Headers, mediaParam, []byte(value))
				}
			}
		}
	}
	writer.Close()
	return *body, writer.FormDataContentType(), nil
}

// formParamWriter writes a form parameter to the multipart writer.
func formParamWriter(
	writer *multipart.Writer,
	fpHeaders http.Header,
	mediaParam map[string]string,
	bytes []byte) error {

	mimeHeader := make(textproto.MIMEHeader)

	contentDisp := mime.FormatMediaType("form-data", mediaParam)
	mimeHeader.Set("Content-Disposition", contentDisp)

	if contentType := fpHeaders.Get("Content-Type"); contentType != "" {
		mimeHeader.Set("Content-Type", contentType)
	}
	part, err := writer.CreatePart(mimeHeader)
	if err != nil {
		return err
	}
	_, err = part.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// structToMap converts a given data structure to a map.
func structToMap2(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	return mapData, err
}

func structToMap(in interface{}) (map[string]interface{}, error) {
	tagName := "json"
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // Non-structural return error
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// Traversing structure fields
	// Specify the tagName value as the key in the map; the field value as the value in the map
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}
