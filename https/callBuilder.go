package https

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
type CallBuilderFactory func(ctx context.Context, httpMethod, path string) CallBuilder
type baseUrlProvider func(server string) string

type CallBuilder interface {
	AppendPath(path string)
	AppendTemplateParam(param string)
	AppendTemplateParams(params interface{})
	BaseUrl(arg string)
	Method(httpMethodName string)
	validateMethod() error
	Accept(acceptHeaderValue string)
	ContentType(contentTypeHeaderValue string)
	Header(name string, value interface{})
	CombineHeaders(headersToMerge map[string]string)
	QueryParam(name string, value interface{})
	validateQueryParams() error
	QueryParams(parameters map[string]interface{})
	FormParam(name string, value interface{})
	validateFormParams() error
	FormData(fields []FormParam)
	validateFormData() error
	Text(body string)
	FileStream(file FileWrapper)
	Json(data interface{})
	validateJson() error
	intercept(interceptor HttpInterceptor)
	InterceptRequest(interceptor func(httpRequest *http.Request) *http.Request)
	toRequest() (*http.Request, error)
	Call() (*HttpContext, error)
	CallAsJson() (*json.Decoder, *http.Response, error)
	CallAsText() (string, *http.Response, error)
	CallAsStream() ([]byte, *http.Response, error)
	Authenticate(requiresAuth bool)
	RequestRetryOption(option RequestRetryOption)
}

type defaultCallBuilder struct {
	ctx                    context.Context
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
	retryOption            RequestRetryOption
	retryConfig            RetryConfiguration
	clientError            error
	jsonData               interface{}
	formFields             []FormParam
	formParams             []FormParam
	queryParams            []FormParam
}

func newDefaultCallBuilder(
	ctx context.Context,
	httpClient HttpClient,
	httpMethod,
	path string,
	baseUrlProvider baseUrlProvider,
	authProvider Authenticator,
	retryConfig RetryConfiguration,
) *defaultCallBuilder {
	cb := defaultCallBuilder{
		httpClient:      httpClient,
		path:            path,
		httpMethod:      httpMethod,
		authProvider:    authProvider,
		baseUrlProvider: baseUrlProvider,
		retryOption:     RequestRetryOption(Default),
		clientError:     nil,
		retryConfig:     retryConfig,
		ctx:             ctx,
	}
	cb.addRetryInterceptor()
	return &cb
}

func (cb *defaultCallBuilder) addAuthentication() {
	if cb.authProvider != nil {
		cb.intercept(cb.authProvider(cb.requiresAuth))
	}
}

func (cb *defaultCallBuilder) Authenticate(requiresAuth bool) {
	cb.requiresAuth = requiresAuth
	if cb.requiresAuth {
		cb.addAuthentication()
	}
}

func (cb *defaultCallBuilder) RequestRetryOption(option RequestRetryOption) {
	cb.retryOption = option
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

func (cb *defaultCallBuilder) Method(httpMethodName string) {
	cb.httpMethod = httpMethodName
}

func (cb *defaultCallBuilder) validateMethod() error {
	if strings.EqualFold(cb.httpMethod, http.MethodGet) {
		cb.httpMethod = http.MethodGet
	} else if strings.EqualFold(cb.httpMethod, http.MethodPut) {
		cb.httpMethod = http.MethodPut
	} else if strings.EqualFold(cb.httpMethod, http.MethodPost) {
		cb.httpMethod = http.MethodPost
	} else if strings.EqualFold(cb.httpMethod, http.MethodPatch) {
		cb.httpMethod = http.MethodPatch
	} else if strings.EqualFold(cb.httpMethod, http.MethodDelete) {
		cb.httpMethod = http.MethodDelete
	} else {
		return internalError{Body: "invalid HTTP method given", FileInfo: "CallBuilder.go/validateMethod"}
	}
	return nil
}

func (cb *defaultCallBuilder) Accept(acceptHeaderValue string) {
	cb.acceptHeaderValue = acceptHeaderValue
	cb.setAcceptIfNotSet(acceptHeaderValue)
}

func (cb *defaultCallBuilder) ContentType(contentTypeHeaderValue string) {
	cb.contentTypeHeaderValue = contentTypeHeaderValue
	cb.setContentTypeIfNotSet(contentTypeHeaderValue)
}

func (cb *defaultCallBuilder) Header(
	name string,
	value interface{},
) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	SetHeaders(cb.headers, strings.ToLower(name), fmt.Sprintf("%v", value))
}

func (cb *defaultCallBuilder) CombineHeaders(headersToMerge map[string]string) {
	MergeHeaders(cb.headers, headersToMerge)
}

func (cb *defaultCallBuilder) QueryParam(
	name string,
	value interface{},
) {
	cb.queryParams = append(cb.queryParams, FormParam{name, value, nil})
}

func (cb *defaultCallBuilder) validateQueryParams() error {
	var err error = nil
	if len(cb.queryParams) != 0 {
		if cb.query == nil {
			cb.query = url.Values{}
		}
		for _, param := range cb.queryParams {
			cb.query, err = prepareFormFields(FormParam{Key: param.Key, Value: param.Value}, cb.query)
			if err != nil {
				return internalError{Body: err.Error(), FileInfo: "CallBuilder.go/validateQueryParams"}
			}
		}
	}
	return nil
}

func (cb *defaultCallBuilder) QueryParams(parameters map[string]interface{}) {
	cb.query = utilities.PrepareQueryParams(cb.query, parameters)
}

func (cb *defaultCallBuilder) FormParam(
	name string,
	value interface{},
) {
	cb.formParams = append(cb.formParams, FormParam{name, value, nil})
}

func (cb *defaultCallBuilder) validateFormParams() error {
	var err error = nil
	if len(cb.formParams) != 0 {
		if cb.form == nil {
			cb.form = url.Values{}
		}
		for _, param := range cb.formParams {
			cb.form, err = prepareFormFields(FormParam{Key: param.Key, Value: param.Value}, cb.form)
			if err != nil {
				return internalError{Body: err.Error(), FileInfo: "CallBuilder.go/validateFormParams"}
			}
			cb.setContentTypeIfNotSet(FORM_URLENCODED_CONTENT_TYPE)
		}
	}
	return nil
}

func (cb *defaultCallBuilder) FormData(fields []FormParam) {
	if fields != nil {
		cb.formFields = fields
	}
}

func (cb *defaultCallBuilder) validateFormData() error {
	var headerVal string
	var err error = nil
	if len(cb.formFields) != 0 {

		cb.formData, headerVal, err = prepareMultipartFields(cb.formFields)
		if err != nil {
			return internalError{Body: err.Error(), FileInfo: "CallBuilder.go/validateFormData"}
		}
		cb.setContentTypeIfNotSet(headerVal)
	}
	return nil
}

func (cb *defaultCallBuilder) Text(body string) {
	cb.body = body
	cb.setContentTypeIfNotSet(TEXT_CONTENT_TYPE)
}

func (cb *defaultCallBuilder) FileStream(file FileWrapper) {
	cb.streamBody = file.File
	cb.setContentTypeIfNotSet("application/octet-stream")
}

func (cb *defaultCallBuilder) Json(data interface{}) {
	cb.jsonData = data
}

func (cb *defaultCallBuilder) validateJson() error {
	if cb.jsonData != nil {
		bytes, err := json.Marshal(cb.jsonData)
		if err != nil {
			return internalError{Body: fmt.Sprintf("Unable to marshal the given data: %v", err.Error()), FileInfo: "CallBuilder.go/validateJson"}
		}
		cb.body = string(bytes)
		cb.setContentTypeIfNotSet(JSON_CONTENT_TYPE)
	}
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

func (cb *defaultCallBuilder) setAcceptIfNotSet(accept string) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	if cb.headers[ACCEPT_HEADER] == "" {
		cb.headers[ACCEPT_HEADER] = accept
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
	var err error
	request := http.Request{}

	err = cb.validateMethod()
	if err != nil {
		return &request, err
	} else {
		request.Method = cb.httpMethod
	}

	url, err := url.Parse(mergePath(cb.baseUrlProvider(cb.baseUrlArg), cb.path))
	if err != nil {
		return &request, err
	}

	err = cb.validateQueryParams()
	if err != nil {
		return &request, err
	} else {
		if len(cb.query) > 0 {
			url.RawQuery = encodeSpace(cb.query.Encode())
		}
	}

	request.URL = url

	request.Header = make(http.Header)

	err = cb.validateJson()
	if err != nil {
		return &request, err
	} else {
		if strings.TrimSpace(cb.body) != "" {
			request.Body = io.NopCloser(bytes.NewBuffer([]byte(cb.body)))
			defer request.Body.Close()
		}
	}

	err = cb.validateFormData()
	if err != nil {
		return &request, err
	} else {
		if cb.formData.Bytes() != nil {
			request.Body = io.NopCloser(&cb.formData)
		}
	}

	err = cb.validateFormParams()
	if err != nil {
		return &request, err
	} else {
		if len(cb.form) > 0 {
			request.Form = cb.form
			replaced := encodeSpace(cb.form.Encode())
			request.Body = io.NopCloser(bytes.NewBuffer([]byte(replaced)))
		}
	}

	if cb.streamBody != nil {
		request.Body = io.NopCloser(bytes.NewBuffer(cb.streamBody))
	}

	for key, val := range cb.headers {
		if request.Header.Get(key) != "" {
			continue
		} else {
			request.Header.Add(key, val)
		}
	}

	return request.WithContext(cb.ctx), err
}

func (cb *defaultCallBuilder) Call() (*HttpContext, error) {
	f := func(request *http.Request) HttpContext {
		client := cb.httpClient
		response, err := client.Execute(request)
		cb.clientError = err
		return HttpContext{
			Request:  request,
			Response: response,
		}
	}

	pipeline := CallHttpInterceptors(cb.interceptors, f)
	request, err := cb.toRequest()
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, nil, err
	}

	if result.Response != nil {
		if result.Response.Body == http.NoBody {
			err = fmt.Errorf("response body empty")
		}

		return json.NewDecoder(result.Response.Body), result.Response, err
	}
	if cb.clientError != nil {
		return nil, nil, cb.clientError
	}
	return nil, nil, err
}

func (cb *defaultCallBuilder) CallAsText() (string, *http.Response, error) {
	result, err := cb.Call()
	if err != nil {
		return "", nil, err
	}
	if result.Response != nil {
		if result.Response.Body == http.NoBody {
			return "", result.Response, fmt.Errorf("response body empty")
		}

		body, err := io.ReadAll(result.Response.Body)
		if err != nil {
			return "", result.Response, fmt.Errorf("Error reading Response body: %v", err.Error())
		}

		return string(body), result.Response, err
	}
	if cb.clientError != nil {
		return "", nil, cb.clientError
	}
	return "", nil, err
}

func (cb *defaultCallBuilder) CallAsStream() ([]byte, *http.Response, error) {
	result, err := cb.Call()
	if err != nil {
		return nil, nil, err
	}

	if result.Response != nil {
		if result.Response.Body == http.NoBody {
			return nil, result.Response, fmt.Errorf("response body empty")
		}

		bytes, err := io.ReadAll(result.Response.Body)
		if err != nil {
			return nil, result.Response, fmt.Errorf("Error reading Response body: %v", err.Error())
		}

		return bytes, result.Response, err
	}
	if cb.clientError != nil {
		return nil, nil, cb.clientError
	}
	return nil, nil, err
}

func (cb *defaultCallBuilder) addRetryInterceptor() {
	cb.intercept(
		func(
			req *http.Request,
			next HttpCallExecutor,
		) HttpContext {
			var context HttpContext
			allowedWaitTime := cb.retryConfig.maximumRetryWaitTime
			if allowedWaitTime == 0 {
				allowedWaitTime = 1<<63 - 1
			}
			shouldRetry := cb.retryConfig.ShouldRetry(cb.retryOption, req.Method)
			retryCount := 0
			var waitTime time.Duration

			for ok := true; ok; ok = waitTime > 0 {
				select {
				case <-req.Context().Done():
					cb.clientError = fmt.Errorf("Request cancelled: %v", req.Context().Err())
					return HttpContext{Request: req}
				default:
				}

				context = next(req)
				if retryCount > 0 {
					allowedWaitTime -= waitTime
				}
				if shouldRetry {
					waitTime = cb.retryConfig.GetRetryWaitTime(
						allowedWaitTime,
						int64(retryCount),
						context.Response,
						cb.clientError)
					time.Sleep(time.Duration(waitTime) * time.Second)
					retryCount++
				}
			}
			return context
		})
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
	retryConfig RetryConfiguration,
) CallBuilderFactory {
	return func(
		ctx context.Context,
		httpMethod,
		path string,
	) CallBuilder {
		return newDefaultCallBuilder(
			ctx,
			httpClient,
			httpMethod,
			path,
			baseUrlProvider,
			auth,
			retryConfig,
		)
	}
}
