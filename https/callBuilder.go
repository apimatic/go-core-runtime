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

// Constants for commonly used HTTP headers and content types.
const CONTENT_TYPE_HEADER = "content-type"
const ACCEPT_HEADER = "accept"
const CONTENT_LENGTH_HEADER = "content-length"
const AUTHORIZATION_HEADER = "authorization"
const FORM_URLENCODED_CONTENT_TYPE = "application/x-www-form-urlencoded"
const JSON_CONTENT_TYPE = "application/json"
const TEXT_CONTENT_TYPE = "text/plain; charset=utf-8"
const XML_CONTENT_TYPE = "application/xml"
const MULTIPART_CONTENT_TYPE = "multipart/form-data"

// CallBuilderFactory is a function type used to create CallBuilder instances for making API calls.
type CallBuilderFactory func(ctx context.Context, httpMethod, path string) CallBuilder

type ErrorBuilder func(code int) error

// baseUrlProvider is a function type used to provide the base URL for API calls based on the server name.
type baseUrlProvider func(server string) string

// CallBuilder is an interface that defines methods for building and executing HTTP requests for API calls.
type CallBuilder interface {
	AppendPath(path string)
	AppendTemplateParam(param string)
	AppendTemplateParams(params interface{})
	AppendErrors(errors map[string]ErrorBuilder)
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
	FormData(fields FormParams)
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
	Authenticate(authGroup AuthGroup)
	RequestRetryOption(option RequestRetryOption)
}

// defaultCallBuilder is a struct that implements the CallBuilder interface for making API calls.
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
	authProvider           map[string]AuthInterface
	retryOption            RequestRetryOption
	retryConfig            RetryConfiguration
	clientError            error
	jsonData               interface{}
	formFields             FormParams
	formParams             FormParams
	queryParams            FormParams
	errors                 map[string]ErrorBuilder
}

// newDefaultCallBuilder creates a new instance of defaultCallBuilder, which implements the CallBuilder interface.
func newDefaultCallBuilder(
	ctx context.Context,
	httpClient HttpClient,
	httpMethod,
	path string,
	baseUrlProvider baseUrlProvider,
	authProvider map[string]AuthInterface,
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

// Authenticate sets the authentication requirement for the API call.
// If a valid auth is given, it adds the respective authentication interceptor to the CallBuilder.
func (cb *defaultCallBuilder) Authenticate(authGroup AuthGroup) {

	authGroup.validate(cb.authProvider)

	if authGroup.errMessage != "" {
		cb.clientError = internalError{
			Type:     AUTHENTICATION_ERROR,
			Body:     authGroup.errMessage,
			FileInfo: "callBuilder.go/Authenticate",
		}
		return
	}

	for _, authI := range authGroup.validAuthInterfaces {
		cb.intercept(authI.Authenticator())
	}

}

// RequestRetryOption sets the retry option for the API call.
// It allows users to configure the retry behavior for the request.
func (cb *defaultCallBuilder) RequestRetryOption(option RequestRetryOption) {
	cb.retryOption = option
}

// AppendPath appends the provided path to the existing path in the CallBuilder.
func (cb *defaultCallBuilder) AppendPath(path string) {
	if cb.path != "" {
		cb.path = sanitizePath(mergePath(cb.path, path))
	} else {
		cb.path = sanitizePath(path)
	}
}

// AppendTemplateParam appends the provided parameter to the existing path in the CallBuilder as a URL template parameter.
func (cb *defaultCallBuilder) AppendTemplateParam(param string) {
	if strings.Contains(cb.path, "%s") {
		cb.path = fmt.Sprintf(cb.path, "/"+url.QueryEscape(param))
	} else {
		cb.AppendPath(url.QueryEscape(param))
	}
}

// AppendTemplateParams appends the provided parameters to the existing path in the CallBuilder as URL template parameters.
// It accepts a slice of strings or a slice of integers as the params argument.
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

func (cb *defaultCallBuilder) AppendErrors(errors map[string]ErrorBuilder) {
	if cb.errors == nil {
		cb.errors = make(map[string]ErrorBuilder)
	}
	for key, err := range errors {
		cb.errors[key] = err
	}
}

// BaseUrl sets the base URL for the API call.
// It takes the server name as an argument and uses the baseUrlProvider function to get the actual base URL.
func (cb *defaultCallBuilder) BaseUrl(server string) {
	cb.baseUrlArg = server
}

// Method sets the HTTP method for the API call.
// It validates the provided HTTP method and stores it in the CallBuilder.
func (cb *defaultCallBuilder) Method(httpMethodName string) {
	cb.httpMethod = httpMethodName
}

// validateMethod validates the HTTP method set in the CallBuilder.
// If the method is not one of the standard HTTP methods, it returns an error.
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

// Accept sets the "Accept" header for the API call.
// It takes the acceptHeaderValue as an argument and sets it as the value for the "Accept" header in the CallBuilder.
func (cb *defaultCallBuilder) Accept(acceptHeaderValue string) {
	cb.acceptHeaderValue = acceptHeaderValue
	cb.setAcceptIfNotSet(acceptHeaderValue)
}

// ContentType sets the "Content-Type" header for the API call.
// It takes the contentTypeHeaderValue as an argument and sets it as the value for the "Content-Type" header in the CallBuilder.
func (cb *defaultCallBuilder) ContentType(contentTypeHeaderValue string) {
	cb.contentTypeHeaderValue = contentTypeHeaderValue
	cb.setContentTypeIfNotSet(contentTypeHeaderValue)
}

// Header sets a custom header for the API call.
// It takes the name of the header and the value of the header as arguments.
func (cb *defaultCallBuilder) Header(
	name string,
	value interface{},
) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	SetHeaders(cb.headers, strings.ToLower(name), fmt.Sprintf("%v", value))
}

// CombineHeaders combines the provided headers with the existing headers in the CallBuilder.
func (cb *defaultCallBuilder) CombineHeaders(headersToMerge map[string]string) {
	MergeHeaders(cb.headers, headersToMerge)
}

// QueryParam adds a query parameter to the API call.
// It takes the name and value of the query parameter as arguments.
func (cb *defaultCallBuilder) QueryParam(
	name string,
	value interface{},
) {
	cb.queryParams.Add(FormParam{name, value, nil})
}

// validateQueryParams validates the query parameters in the CallBuilder.
func (cb *defaultCallBuilder) validateQueryParams() error {
	if len(cb.queryParams) != 0 {
		if cb.query == nil {
			cb.query = url.Values{}
		}
		err := cb.queryParams.prepareFormFields(cb.query)
		if err != nil {
			return internalError{Body: err.Error(), FileInfo: "CallBuilder.go/validateQueryParams"}
		}
	}
	return nil
}

// QueryParams sets multiple query parameters for the API call.
// It takes a map of string keys and interface{} values representing the query parameters.
func (cb *defaultCallBuilder) QueryParams(parameters map[string]interface{}) {
	cb.query = utilities.PrepareQueryParams(cb.query, parameters)
}

// FormParam adds a form parameter to the API call.
// It takes the name and value of the form parameter as arguments.
func (cb *defaultCallBuilder) FormParam(
	name string,
	value interface{},
) {
	cb.formParams.Add(FormParam{name, value, nil})
}

// validateFormParams validates the form parameters in the CallBuilder.
// Additionally, it sets the "Content-Type" header to "application/x-www-form-urlencoded" if not already set.
func (cb *defaultCallBuilder) validateFormParams() error {
	if len(cb.formParams) != 0 {
		if cb.form == nil {
			cb.form = url.Values{}
		}
		err := cb.formParams.prepareFormFields(cb.form)
		if err != nil {
			return internalError{Body: err.Error(), FileInfo: "CallBuilder.go/validateFormParams"}
		}
		cb.setContentTypeIfNotSet(FORM_URLENCODED_CONTENT_TYPE)
	}
	return nil
}

// FormData sets form fields for the API call.
// It takes a slice of FormParam representing the form fields.
func (cb *defaultCallBuilder) FormData(fields FormParams) {
	if fields != nil {
		cb.formFields = fields
	}
}

// validateFormData validates the form fields in the CallBuilder.
// Additionally, it sets the "Content-Type" header to the appropriate value for multipart form data if not already set.
func (cb *defaultCallBuilder) validateFormData() error {
	var headerVal string
	var err error = nil
	if len(cb.formFields) != 0 {
		cb.formData, headerVal, err = cb.formFields.prepareMultipartFields()
		if err != nil {
			return internalError{Body: err.Error(), FileInfo: "CallBuilder.go/validateFormData"}
		}
		cb.setContentTypeIfNotSet(headerVal)
	}
	return nil
}

// Text sets the request body for the API call as plain text.
// It takes the body string as an argument.
// Additionally, it sets the "Content-Type" header to "text/plain; charset=utf-8" if not already set.
func (cb *defaultCallBuilder) Text(body string) {
	cb.body = body
	cb.setContentTypeIfNotSet(TEXT_CONTENT_TYPE)
}

// FileStream sets the request body for the API call as a file stream.
// It takes a FileWrapper struct containing the file content as an argument.
// Additionally, it sets the "Content-Type" header to "application/octet-stream" if not already set.
func (cb *defaultCallBuilder) FileStream(file FileWrapper) {
	cb.streamBody = file.File
	cb.setContentTypeIfNotSet("application/octet-stream")
}

// Json sets the request body for the API call as JSON.
func (cb *defaultCallBuilder) Json(data interface{}) {
	cb.jsonData = data
}

// validateJson validates the JSON data in the CallBuilder.
// It marshals the JSON data into a byte array and stores it as the request body.
// Additionally, it sets the "Content-Type" header to "application/json" if not already set.
// If there is an error during marshaling, it returns an internalError.
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

// setContentTypeIfNotSet sets the "Content-Type" header if it is not already set in the CallBuilder.
// It takes the contentType as an argument and sets it as the value for the "Content-Type" header.
// If the headers map is nil, it initializes it before setting the header.
func (cb *defaultCallBuilder) setContentTypeIfNotSet(contentType string) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	if cb.headers[CONTENT_TYPE_HEADER] == "" {
		cb.headers[CONTENT_TYPE_HEADER] = contentType
	}
}

// setAcceptIfNotSet sets the "Accept" header if it is not already set in the CallBuilder.
// It takes the accept header value as an argument and sets it as the value for the "Accept" header.
// If the headers map is nil, it initializes it before setting the header.
func (cb *defaultCallBuilder) setAcceptIfNotSet(accept string) {
	if cb.headers == nil {
		cb.headers = make(map[string]string)
	}
	if cb.headers[ACCEPT_HEADER] == "" {
		cb.headers[ACCEPT_HEADER] = accept
	}
}

// intercept adds the provided HTTP interceptor to the list of interceptors in the CallBuilder.
// This allows users to add custom HTTP interceptors for modifying the request and response.
func (cb *defaultCallBuilder) intercept(interceptor HttpInterceptor) {
	cb.interceptors = append(cb.interceptors, interceptor)
}

// InterceptRequest adds an interceptor function for modifying the request before sending.
// The interceptor function takes the original http.Request as input and returns a modified http.Request.
// The modified request is used for making the API call.
// Use this method to customize the request headers, query parameters, or other attributes.
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

// toRequest converts the CallBuilder configuration into an http.Request object.
// It prepares the request by setting the HTTP method, URL, headers, and request body.
// If there are any validation errors, it returns an error along with an empty request.
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

// Call executes the API call and returns the HttpContext that contains the request and response.
// It iterates through the interceptors to execute them in sequence before making the API call.
func (cb *defaultCallBuilder) Call() (*HttpContext, error) {
	// return any client errors found before executing the call
	if cb.clientError != nil {
		return nil, cb.clientError
	}

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

// CallAsJson executes the API call and returns a JSON decoder and the HTTP response.
// It sets the "Accept" header to "application/json" and calls the Call method to make the API call.
// The JSON decoder allows users to parse the response body as JSON data directly.
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

// CallAsText executes the API call and returns the response body as a string and the HTTP response.
// It calls the Call method to make the API call and reads the response body as a byte array.
// The byte array is then converted to a string and returned as the response body.
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

// CallAsStream executes the API call and returns the response body as a byte array and the HTTP response.
// It calls the Call method to make the API call and reads the response body as a byte array.
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

// addRetryInterceptor adds a retry interceptor to the call builder. This interceptor will handle retrying the API call
// based on the provided retry configuration. It checks if the call should be retried, and if so, waits for the specified
// amount of time before making the next retry attempt. The retry logic continues until either the maximum retry wait time
// is exceeded or the request is successfully executed.
// The result of the last API call is returned in the HttpContext through the interceptor.
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
					cb.clientError = fmt.Errorf("request cancelled: %v", req.Context().Err())
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

// mergePath combines two URL paths to create a valid URL path.
func mergePath(left, right string) string {
	if right == "" {
		return left
	}

	if left[len(left)-1] == '/' && right[0] == '/' {
		return left + strings.Replace(right, "/", "", 1)
	} else if left[len(left)-1] == '/' || right[0] == '/' {
		return left + right
	} else {
		return left + "/" + right
	}
}

// sanitizePath removes any duplicate slashes in the given path to create a valid URL path.
func sanitizePath(path string) string {
	return strings.Replace(path, "//", "/", -1)
}

// encodeSpace replaces all occurrences of the plus sign "+" with "%20" in the given string.
// This function is used to encode spaces in query parameters.
func encodeSpace(str string) string {
	return strings.ReplaceAll(str, "+", "%20")
}

// CreateCallBuilderFactory creates a new CallBuilderFactory function which
// creates a new CallBuilder using the provided inputs
func CreateCallBuilderFactory(
	baseUrlProvider baseUrlProvider,
	auth map[string]AuthInterface,
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
