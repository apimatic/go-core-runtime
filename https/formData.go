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
func (fp *FormParams) prepareFormFields(form url.Values, option ArraySerializationOption) error {
	if form == nil {
		form = url.Values{}
	}
	for _, param := range *fp {
		paramsMap, err := toMap(param.Key, param.Value, option)
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
func (fp *FormParams) prepareMultipartFields(option ArraySerializationOption) (bytes.Buffer, string, error) {
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
			paramsMap, err := toMap(field.Key, field.Value, option)
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

func toMap(keyPrefix string, param any, option ArraySerializationOption) (map[string][]string, error) {
	if param == nil {
		return map[string][]string{}, nil
	}

	valueKind := reflect.TypeOf(param).Kind()
	switch valueKind {
	case reflect.Struct, reflect.Ptr:
		innerMap, err := structToMap(param)
		if err != nil {
			return nil, err
		}
		innerFlatMap, err := toMap(keyPrefix, innerMap, option)
		if err != nil {
			return nil, err
		}
		return innerFlatMap, nil
	case reflect.Map:
		iter := reflect.ValueOf(param).MapRange()
		result := make(map[string][]string)
		for iter.Next() {
			innerKey := ArraySerializationOption(Indexed).joinKey(keyPrefix, iter.Key())
			innerValue := iter.Value().Interface()
			innerFlatMap, err := toMap(innerKey, innerValue, option)
			if err != nil {
				return nil, err
			}
			option.appendMap(result, innerFlatMap)
		}
		return result, nil
	case reflect.Slice:
		reflectValue := reflect.ValueOf(param)
		result := make(map[string][]string)
		for i := 0; i < reflectValue.Len(); i++ {
			var innerStruct = reflectValue.Index(i).Interface()
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
	default:
		return map[string][]string{keyPrefix: {getDefaultValue(param)}}, nil
	}
}

// structToMap converts a given data structure to a map.
func structToMap(data any) (map[string]any, error) {
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		data = toStructPtr(data)
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	return mapData, err
}

// Return a pointer to the supplied struct via interface{}
func toStructPtr(obj any) any {
	// Create a new instance of the underlying type
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	// NOTE: `vp.Elem().Set(reflect.ValueOf(&obj).Elem())` does not work
	// Return a `Cat` pointer to obj -- i.e. &obj.(*Cat)
	return vp.Interface()
}

func getDefaultValue(in any) string {
	switch in.(type) {
	case string:
		return in.(string)
	default:
		bytes, err := json.Marshal(in)
		if err == nil {
			return string(bytes)
		} else {
			return fmt.Sprintf("%v", in)
		}
	}
}
