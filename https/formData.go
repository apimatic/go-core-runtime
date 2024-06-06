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

func (fp *formParam) clone(key string, value any) formParam {
	return formParam{
		key:                      key,
		value:                    value,
		headers:                  fp.headers,
		arraySerializationOption: fp.arraySerializationOption,
	}
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
			paramsMap, err := field.toMap()
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

func (fp *formParam) IsMultipart() bool {
	contentType := fp.headers.Get(CONTENT_TYPE_HEADER)
	if contentType != "" {
		return contentType != FORM_URLENCODED_CONTENT_TYPE
	}
	return false
}

func (fp *formParam) toMap() (map[string][]string, error) {
	if fp.value == nil {
		return map[string][]string{}, nil
	}

	if (fp.IsMultipart()){
		return fp.processDefault()
	}
	
	switch reflect.TypeOf(fp.value).Kind() {
	case reflect.Ptr:
		return fp.processStructAndPtr()
	case reflect.Struct:
		innerfp := fp.clone(fp.key, toStructPtr(fp.value))
		return innerfp.processStructAndPtr()
	case reflect.Map:
		return fp.processMap()
	case reflect.Slice:
		return fp.processSlice()
	default:
		return fp.processDefault()
	}
}

func (fp *formParam) processStructAndPtr() (map[string][]string, error) {
	innerData, err := structToAny(fp.value)
	if err != nil {
		return nil, err
	}

	innerfp := fp.clone(fp.key, innerData)
	return innerfp.toMap()
}

func (fp *formParam) processMap() (map[string][]string, error) {
	iter := reflect.ValueOf(fp.value).MapRange()
	result := make(map[string][]string)
	for iter.Next() {
		innerKey := fp.arraySerializationOption.joinKey(fp.key, iter.Key().Interface())
		innerValue := iter.Value().Interface()
		innerfp := fp.clone(innerKey, innerValue)
		innerFlatMap, err := innerfp.toMap()
		if err != nil {
			return nil, err
		}
		fp.arraySerializationOption.appendMap(result, innerFlatMap)
	}
	return result, nil
}

func (fp *formParam) processSlice() (map[string][]string, error) {
	reflectValue := reflect.ValueOf(fp.value)
	result := make(map[string][]string)
	for i := 0; i < reflectValue.Len(); i++ {
		innerElem := reflectValue.Index(i)
		var indexStr any
		switch innerElem.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
			indexStr = nil
		default:
			indexStr = fmt.Sprintf("%v", i)
		}
		innerKey := fp.arraySerializationOption.joinKey(fp.key, indexStr)
		innerfp := fp.clone(innerKey, innerElem.Interface())
		innerFlatMap, err := innerfp.toMap()
		if err != nil {
			return result, err
		}
		fp.arraySerializationOption.appendMap(result, innerFlatMap)
	}
	return result, nil
}

func (fp *formParam) processDefault() (map[string][]string, error) {
	var defaultValue string
	switch reflect.TypeOf(fp.value).Kind() {
	case reflect.String:
		defaultValue = fmt.Sprintf("%v", fp.value)
	default:
		dataBytes, err := json.Marshal(fp.value)
		if err == nil {
			defaultValue = string(dataBytes)
		} else {
			defaultValue = fmt.Sprintf("%v", fp.value)
		}
	}
	return map[string][]string{fp.key: {defaultValue}}, nil
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
