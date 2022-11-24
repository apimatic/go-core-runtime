package https

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/apimatic/go-core-runtime/utilities"
)

const CONTENT_TYPE_HEADER = "content-type"
const ACCEPT_HEADER = "accept"
const CONTENT_LENGTH_HEADER = "content-length"
const AUTHORIZATION_HEADER = "authorization"
const FORM_URLENCODED_CONTENT_TYPE = "application/x-www-form-urlencoded"
const JSON_CONTENT_TYPE = "application/json"
const TEXT_CONTENT_TYPE = "text/plain; charset=utf-8"
const XML_CONTENT_TYPE = "application/xml"
const MULTIPART_CONTENT_TYPE = "multipart/form-data"

type Authenticator func(bool) HttpInterceptor
type CallBuilderFactory func(httpMethod, path string) CallBuilder
type baseUrlProvider func(server string) string

type CallBuilder interface {
	AppendPath(path string)
	AppendTemplateParam(param string)
	AppendTemplateParams(params interface{})
	BaseUrl(arg string)
	Method(httpMethodName string) error
	Accept(acceptHeaderValue string)
	ContentType(contentTypeHeaderValue string)
	Header(name string, value interface{})
	CombineHeaders(headersToMerge map[string]string)
	QueryParam(name string, value interface{}) error
	QueryParams(parameters map[string]interface{})
	FormParam(name string, value interface{}) error
	FormData(fields map[string]interface{}) error
	Text(body string)
	FileStream(file FileWrapper)
	Json(data interface{}) error
	intercept(interceptor HttpInterceptor)
	InterceptRequest(interceptor func(httpRequest *http.Request) *http.Request)
	toRequest() (*http.Request, error)
	Call() (*HttpContext, error)
	CallAsJson() (*json.Decoder, *http.Response, error)
	CallAsText() (string, *http.Response, error)
	CallAsStream() ([]byte, *http.Response, error)
	Authenticate(requiresAuth bool)
}

type defaultCallBuilder struct {
	path                   string
	baseUrlArg             string
	baseUrlProvider        baseUrlProvider
	httpMethod             string
	acceptHeaderValue      string
	contentTypeHeaderValue string
	headers                map[string]string
	query                  url.Values
	form                   url.Values
	formData               bytes.Buffer
	body                   string
	streamBody             []byte
	httpClient             HttpClient
	interceptors           []HttpInterceptor
	requiresAuth           bool
	authProvider           Authenticator
}

func newDefaultCallBuilder(
	httpClient HttpClient,
	httpMethod,
	path string,
	baseUrlProvider baseUrlProvider,
	authProvider Authenticator,
) *defaultCallBuilder {
	cb := defaultCallBuilder{
		httpClient:      httpClient,
		path:            path,
		httpMethod:      httpMethod,
		authProvider:    authProvider,
		baseUrlProvider: baseUrlProvider,
	}
	return &cb
}

func (cb *defaultCallBuilder) addAuthentication() {
	if cb.authProvider != nil {
		cb.intercept(cb.authProvider(cb.requiresAuth))
	}
}

func (cb *defaultCallBuilder) Authenticate(requiresAuth bool) {
	cb.requiresAuth = requiresAuth
	if cb.requiresAuth == true {
		cb.addAuthentication()
	}
}

func (cb *defaultCallBuilder) AppendPath(path string) {
	if cb.path != "" {
		cb.path = sanitizePath(mergePath(cb.path, path))
	} else {
		cb.path = sanitizePath(path)
	}
}

func (cb *defaultCallBuilder) AppendTemplateParam(param string) {
	if strings.Contains(cb.path, "%s") {
		cb.path = fmt.Sprintf(cb.path, "/"+url.QueryEscape(param))
	} else {
		cb.AppendPath(url.QueryEscape(param))
	}
}

func (cb *defaultCallBuilder) AppendTemplateParams(params interface{}) {
	switch x := params.(type) {
	case []string:
		for _, param := range x {
			cb.AppendTemplateParam(param)
		}
	case []int:
		for _, param := range x {
			cb.AppendTemplateParam(strconv.Itoa(int(param)))
		}
	}
}

func (cb *defaultCallBuilder) BaseUrl(server string) {
	cb.baseUrlArg = server
}

func (cb *defaultCallBuilder) Method(httpMethodName string) error {
	if strings.EqualFold(httpMethodName, http.MethodGet) {
		cb.httpMethod = http.MethodGet
	} else if strings.EqualFold(httpMethodName, http.MethodPut) {
		cb.httpMethod = http.MethodPut
	} else if strings.EqualFold(httpMethodName, http.MethodPost) {
		cb.httpMethod = http.MethodPost
	} else if strings.EqualFold(httpMethodName, http.MethodPatch) {
		cb.httpMethod = http.MethodPatch
	} else if strings.EqualFold(httpMethodName, http.MethodDelete) {
		cb.httpMethod = http.MethodDelete
	} else {
		return fmt.Errorf("Invalid HTTP method given!")
	}
	return nil
}

func (cb *defaultCallBuilder) Accept(acceptHeaderValue string) {
	cb.acceptHeaderValue = acceptHeaderValue
}

func (cb *defaultCallBuilder) ContentType(contentTypeHeaderValue string) {
	cb.contentTypeHeaderValue = contentTypeHeaderValue
}

func (cb *defaultCallBuilder) Header(
	name string,
	value interface{},
) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	SetHeaders(cb.headers, name, fmt.Sprintf("%v", value))
}

func (cb *defaultCallBuilder) CombineHeaders(headersToMerge map[string]string) {
	MergeHeaders(cb.headers, headersToMerge)
}

func (cb *defaultCallBuilder) QueryParam(
	name string,
	value interface{},
) error {
	var err error = nil
	if cb.query == nil {
		cb.query = url.Values{}
	}
	cb.query, err = PrepareFormFields(name, value, cb.query)
	return err
}

func (cb *defaultCallBuilder) QueryParams(parameters map[string]interface{}) {
	cb.query = utilities.PrepareQueryParams(cb.query, parameters)
}

func (cb *defaultCallBuilder) FormParam(
	name string,
	value interface{},
) error {
	var err error = nil
	if cb.form == nil {
		cb.form = url.Values{}
	}
	cb.form, err = PrepareFormFields(name, value, cb.form)
	cb.setContentTypeIfNotSet(FORM_URLENCODED_CONTENT_TYPE)
	return err
}

func (cb *defaultCallBuilder) FormData(fields map[string]interface{}) error {
	var headerVal string
	var err error = nil
	cb.formData, headerVal, err = PrepareMultipartFields(fields)
	cb.setContentTypeIfNotSet(headerVal)
	return err
}

func (cb *defaultCallBuilder) Text(body string) {
	cb.body = body
	cb.setContentTypeIfNotSet(TEXT_CONTENT_TYPE)
}

func (cb *defaultCallBuilder) FileStream(file FileWrapper) {
	cb.streamBody = file.File
	if cb.contentTypeHeaderValue != "" {
		cb.Header(CONTENT_TYPE_HEADER, cb.contentTypeHeaderValue)
	} else {
		cb.Header(CONTENT_TYPE_HEADER, "application/octet-stream")
	}
}

func (cb *defaultCallBuilder) Json(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Unable to marshal the given data: %v", err)
	}
	cb.body = string(bytes)
	cb.setContentTypeIfNotSet(JSON_CONTENT_TYPE)
	return nil
}

func (cb *defaultCallBuilder) setContentTypeIfNotSet(contentType string) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	if cb.headers[CONTENT_TYPE_HEADER] == "" {
		cb.headers[CONTENT_TYPE_HEADER] = contentType
	}
}

func (cb *defaultCallBuilder) intercept(interceptor HttpInterceptor) {
	cb.interceptors = append(cb.interceptors, interceptor)
}

func (cb *defaultCallBuilder) InterceptRequest(
	interceptor func(httpRequest *http.Request) *http.Request,
) {
	cb.intercept(
		func(
			req *http.Request,
			next HttpCallExecutor,
		) HttpContext {
			return next(interceptor(req))
		})
}

func (cb *defaultCallBuilder) toRequest() (*http.Request, error) {
	request := http.Request{
		Method: cb.httpMethod,
	}

	url, err := url.Parse(mergePath(cb.baseUrlProvider(cb.baseUrlArg), cb.path))
	if len(cb.query) > 0 {
		url.RawQuery = encodeSpace(cb.query.Encode())
	}
	request.URL = url

	request.Header = make(http.Header)
	if strings.TrimSpace(cb.acceptHeaderValue) != "" {
		request.Header.Add(ACCEPT_HEADER, cb.acceptHeaderValue)
	}

	if strings.TrimSpace(cb.contentTypeHeaderValue) != "" {
		request.Header.Add(CONTENT_TYPE_HEADER, cb.contentTypeHeaderValue)
	}

	for key, val := range cb.headers {
		if request.Header.Get(key) != "" {
			continue
		} else {
			request.Header.Add(key, val)
		}
	}

	if strings.TrimSpace(cb.body) != "" {
		request.Body = io.NopCloser(bytes.NewBuffer([]byte(cb.body)))
		defer request.Body.Close()
	} else if cb.formData.Bytes() != nil {
		request.Body = io.NopCloser(&cb.formData)
	} else if len(cb.form) > 0 {
		request.Form = cb.form
		replaced := encodeSpace(cb.form.Encode())
		request.Body = io.NopCloser(bytes.NewBuffer([]byte(replaced)))
	} else if cb.streamBody != nil {
		request.Body = io.NopCloser(bytes.NewBuffer(cb.streamBody))
	}

	return &request, err
}

func (cb *defaultCallBuilder) Call() (*HttpContext, error) {
	f := func(request *http.Request) HttpContext {
		client := cb.httpClient
		response, _ := client.Execute(request)
		return HttpContext{
			Request:  request,
			Response: response,
		}
	}

	pipeline := CallHttpInterceptors(cb.interceptors, f)
	request, err := cb.toRequest()
	context := pipeline(request)
	return &context, err
}

func (cb *defaultCallBuilder) CallAsJson() (*json.Decoder, *http.Response, error) {
	f := func(request *http.Request) *http.Request {
		request.Header.Set(ACCEPT_HEADER, JSON_CONTENT_TYPE)
		return request
	}

	cb.InterceptRequest(f)
	result, err := cb.Call()
	if result.Response.Body == http.NoBody {
		err = fmt.Errorf("Response body empty!")
	}

	return json.NewDecoder(result.Response.Body), result.Response, err
}

func (cb *defaultCallBuilder) CallAsText() (string, *http.Response, error) {
	result, err := cb.Call()
	if result.Response.Body == http.NoBody {
		err = fmt.Errorf("Response body empty!")
	}

	body, err := ioutil.ReadAll(result.Response.Body)
	if err != nil {
		err = fmt.Errorf("Error reading Response body: %v", err)
	}

	return string(body), result.Response, err
}

func (cb *defaultCallBuilder) CallAsStream() ([]byte, *http.Response, error) {
	result, err := cb.Call()
	if result.Response.Body == http.NoBody {
		err = fmt.Errorf("Response body empty!")
	}

	bytes, err := ioutil.ReadAll(result.Response.Body)
	if err != nil {
		err = fmt.Errorf("Error reading Response body: %v", err)
	}

	return bytes, result.Response, err
}

func mergePath(left, right string) string {
	if right == "" {
		return left
	}

	if left[len(left)-1] == '/' && right[0] == '/' {
		return left + strings.Replace(right, "/", "", -1)
	} else if left[len(left)-1] == '/' || right[0] == '/' {
		return left + right
	} else {
		return left + "/" + right
	}
}

func sanitizePath(path string) string {
	return strings.Replace(path, "//", "/", -1)
}

func encodeSpace(str string) string {
	return strings.ReplaceAll(str, "+", "%20")
}

func CreateCallBuilderFactory(
	baseUrlProvider baseUrlProvider,
	auth Authenticator,
	httpClient HttpClient,
) CallBuilderFactory {
	return func(
		httpMethod,
		path string,
	) CallBuilder {
		return newDefaultCallBuilder(
			httpClient,
			httpMethod,
			path,
			baseUrlProvider,
			auth,
		)
	}
}
