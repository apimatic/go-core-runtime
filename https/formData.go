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

func getDefaultValue(in interface{}) string {

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

func appendMap(param map[string][]string, result map[string][]string) {
	for k, v := range param {
		for _, v1 := range v {
			appendMapValue(k, result, v1)
		}
	}
}

func prePareKey(keyPrefix string, key string) string {
	var innerKey string
	if key == ""  { 		// plain text
		innerKey = fmt.Sprintf("%v", keyPrefix)
	} else if key == "[]" {	// unindexed
		innerKey = fmt.Sprintf("%v[]", keyPrefix)
	} else {				// indexed
		innerKey = fmt.Sprintf("%v[%v]", keyPrefix, key)
	}
	return innerKey
}

func appendMapValue(key string, result map[string][]string, value string) {
	if (len(result[key]) > 0){
		formatter := ","
		result[key][0] = fmt.Sprintf("%v%v%v", result[key][0], formatter, value)
	} else {
		result[key] = append(result[key], value)
	}
	//result[key] = append(result[key], value)
}


// Return a pointer to the supplied struct via interface{}
func toStructPtr(obj interface{}) interface{} {
	// Create a new instance of the underlying type
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	// NOTE: `vp.Elem().Set(reflect.ValueOf(&obj).Elem())` does not work
	// Return a `Cat` pointer to obj -- i.e. &obj.(*Cat)
	return vp.Interface()
}


func toMap3(keyPrefix string, param interface{}) (map[string][]string, error) {
	result := make(map[string][]string)
	valueKind := reflect.TypeOf(param).Kind()

	switch valueKind {
	case reflect.Map:

		iter := reflect.ValueOf(param).MapRange()
		for iter.Next() {
			key := fmt.Sprintf("%v", iter.Key())
			innerKey := prePareKey(keyPrefix, key)
			innerValue := iter.Value()
			innerValueKind := innerValue.Type().Kind()

			var innerStruct any
			if (innerValueKind == reflect.Struct) {
				innerStruct = toStructPtr(innerValue.Interface())
			} else {
				innerStruct =  innerValue.Interface()
			}
			innerFlatMap, err := toMap3(innerKey, innerStruct)
			if err != nil {
				return result, err
			}
			appendMap(innerFlatMap, result)
		}
	case reflect.Struct, reflect.Ptr:
		// Convert Struct and Pointer types into Map.
		innerMap, err := structToMap(toStructPtr(param))
		if err != nil {
			return result, err
		}
		innerFlatMap, err := toMap3(keyPrefix, innerMap)
		if err != nil {
			return result, err
		}
		appendMap(innerFlatMap, result)
	case reflect.Slice:
		reflectValue := reflect.ValueOf(param)
		for index := 0; index < reflectValue.Len(); index++ {
			innerObjType := reflectValue.Index(index).Type().Kind()
			var innerStruct any
			if innerObjType == reflect.Struct {
				innerStruct = reflectValue.Index(index).Addr().Interface()
			} else {
				innerStruct = reflectValue.Index(index).Interface()
			}
			var indexStr string
			switch innerStruct.(type) {
			case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, string:
				indexStr = ""
			default:
				indexStr = fmt.Sprintf("%v", index)
			}
			innerKey := prePareKey(keyPrefix, indexStr)
			innerFlatMap, err := toMap3(innerKey, innerStruct)
			if err != nil {
				return result, err
			}
			appendMap(innerFlatMap, result)
		}
	default:
		appendMapValue(keyPrefix, result, getDefaultValue(param))
	}
	return result, nil
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
		paramsMap, err := toMap3(param.Key, param.Value)
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
			paramsMap, err := toMap3(field.Key, field.Value)
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
func structToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	return mapData, err
}
