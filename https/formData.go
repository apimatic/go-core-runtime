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

type formParam struct {
	key                      string
	value                    any
	headers                  http.Header
	arraySerializationOption ArraySerializationOption
}

type formParams []formParam

// FormParams represents a collection of FormParam objects.
type FormParams []FormParam

// Add appends a FormParam to the FormParams collection.
func (fp *FormParams) Add(formParam FormParam) {
	if formParam.Value != nil {
		*fp = append(*fp, formParam)
	}
}

// Add appends a FormParam to the FormParams collection.
func (fp *formParams) add(formParam formParam) {
	if formParam.value != nil {
		*fp = append(*fp, formParam)
	}
}

// prepareFormFields prepares the form fields from the given FormParams and adds them to the url.Values.
// It processes each FormParam field and encodes the value according to its data type.
func (fp *formParams) prepareFormFields(form url.Values) error {
	if form == nil {
		form = url.Values{}
	}
	for _, param := range *fp {
		paramsMap, err := toMap(param.key, param.value, param.arraySerializationOption)
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
func (fp *formParams) prepareMultipartFields() (bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, field := range *fp {
		switch fieldValue := field.value.(type) {
		case FileWrapper:
			mediaParam := map[string]string{
				"name":     field.key,
				"filename": fieldValue.FileName,
			}
			formParamWriter(writer, field.headers, mediaParam, fieldValue.File)
		default:
			paramsMap, err := toMap(field.key, field.value, field.arraySerializationOption)
			if err != nil {
				return *body, writer.FormDataContentType(), err
			}
			for key, values := range paramsMap {
				mediaParam := map[string]string{"name": key}
				for _, value := range values {
					formParamWriter(writer, field.headers, mediaParam, []byte(value))
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

func toMap(keyPrefix string, param any, option ArraySerializationOption) (map[string][]string, error) {
	if param == nil {
		return map[string][]string{}, nil
	}

	switch paramKind := reflect.TypeOf(param).Kind(); paramKind {
	case reflect.Ptr:
		return processStructAndPtr(keyPrefix, param, option)
	case reflect.Struct:
		return processStructAndPtr(keyPrefix, toStructPtr(param), option)
	case reflect.Map:
		return processMap(keyPrefix, param, option)
	case reflect.Slice:
		return processSlice(keyPrefix, param, option)
	default:
		return processDefault(paramKind, keyPrefix, param)
	}
}

func processStructAndPtr(keyPrefix string, param any, option ArraySerializationOption) (map[string][]string, error) {
	innerData, err := structToAny(param)
	if err != nil { return nil, err }

	switch innerData := innerData.(type) {
	case map[string]any:
		return toMap(keyPrefix, innerData, option)
	default:
		return toMap(keyPrefix, innerData, option)
	} 
}

func processMap(keyPrefix string, param any, option ArraySerializationOption) (map[string][]string, error) {
	iter := reflect.ValueOf(param).MapRange()
	result := make(map[string][]string)
	for iter.Next() {
		innerKey := option.joinKey(keyPrefix, iter.Key().Interface())
		innerValue := iter.Value().Interface()
		innerFlatMap, err := toMap(innerKey, innerValue, option)
		if err != nil {
			return nil, err
		}
		option.appendMap(result, innerFlatMap)
	}
	return result, nil
}

func processSlice(keyPrefix string, param any, option ArraySerializationOption) (map[string][]string, error) {
	reflectValue := reflect.ValueOf(param)
	result := make(map[string][]string)
	for i := 0; i < reflectValue.Len(); i++ {
		innerStruct := reflectValue.Index(i).Interface()
		var indexStr any
		switch innerStruct.(type) {
		case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, string:
			indexStr = nil
		default:
			indexStr = fmt.Sprintf("%v", i)
		}
		innerKey := option.joinKey(keyPrefix, indexStr)
		innerFlatMap, err := toMap(innerKey, innerStruct, option)
		if err != nil {
			return result, err
		}
		option.appendMap(result, innerFlatMap)
	}
	return result, nil
}

func processDefault(paramKind reflect.Kind, keyPrefix string, param any) (map[string][]string, error) {
	var defaultValue string
	switch paramKind {
	case reflect.String:
		defaultValue = fmt.Sprintf("%v", param)
	default:
		dataBytes, err := json.Marshal(param)
		if err == nil {
			defaultValue = string(dataBytes)
		} else {
			defaultValue = fmt.Sprintf("%v", param)
		}
	}
	return map[string][]string{keyPrefix: {defaultValue}}, nil
}

// structToAny converts a given data structure into an any type.
func structToAny(data any) (any, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var innerData any
	err = json.Unmarshal(dataBytes, &innerData)
	return innerData, err
}

// Return a pointer to the supplied struct via any
func toStructPtr(obj any) any {
	// Create a new instance of the underlying type
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	// NOTE: `vp.Elem().Set(reflect.ValueOf(&obj).Elem())` does not work
	// Return a `Cat` pointer to obj -- i.e. &obj.(*Cat)
	return vp.Interface()
}
