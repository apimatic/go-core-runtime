package https

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/assert"
)

var ctx = context.Background()

func GetCallBuilder(ctx context.Context, method, path string, auth map[string]AuthInterface) CallBuilder {
	client := NewHttpClient(NewHttpConfiguration())
	callBuilder := CreateCallBuilderFactory(
		func(server string) string {
			return GetTestingServer().URL
		},
		auth,
		client,
		NewRetryConfiguration(),
		Indexed,
		&ApiLogger{},
	)

	return callBuilder(ctx, method, path)
}

func RequestAuthentication() HttpInterceptor {
	return PassThroughInterceptor
}

func TestAppendPath(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "//response/", nil)
	request.AppendPath("/integer")
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}

	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestAppendMultiplePath(t *testing.T) {
	samplePath := "/number/integer/base64"
	request := GetCallBuilder(ctx, "GET", "//response/", nil)
	request.AppendPath(samplePath)
	_, response, _ := request.CallAsJson()
	responsePath := response.Request.URL.Path
	expected := "/response" + samplePath

	if responsePath != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, responsePath)
	}
}

func TestAppendPathWithMultiSlashes(t *testing.T) {
	samplePath := "//number//integer//base64"
	request := GetCallBuilder(ctx, "GET", "//response/", nil)
	request.AppendPath(samplePath)
	_, response, _ := request.CallAsJson()
	responsePath := response.Request.URL.Path
	expected := "/response" + strings.Replace(samplePath, "//", "/", -1)

	if responsePath != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, responsePath)
	}
}

func TestAppendPathEmptyPath(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.AppendPath("/response/integer")
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}

	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestAppendTemplateParamsStrings(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "/template/%v", nil)
	request.AppendTemplateParams([]string{"abc", "def"})
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}

	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestAppendTemplateParamsIntegers(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "/template/%v", nil)
	request.AppendTemplateParams([]int{1, 2, 3, 4, 5})
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}

	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestBaseUrlValue(t *testing.T) {
	request := GetCallBuilder(ctx, "", "/response/integer", nil)
	request.BaseUrl("https://github.com/apimatic")
	request.Method("GET")
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}
	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestMethodGet(t *testing.T) {
	request := GetCallBuilder(ctx, "", "/response/integer", nil)
	request.Method("GET")
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}
	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestMethodPost(t *testing.T) {
	request := GetCallBuilder(ctx, "", "/form/string", nil)
	request.Method("POST")
	request.FormParam("value", "TestString")
	_, response, err := request.CallAsJson()
	if err != nil {
		t.Errorf("Error in CallAsJson: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", 200, response)
	}
}

func TestMethodPut(t *testing.T) {
	request := GetCallBuilder(ctx, "", "", nil)
	request.Method("PUT")
	result, err := request.toRequest()

	expected := http.MethodPut
	if result.Method != expected || err != nil {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodPatch(t *testing.T) {
	request := GetCallBuilder(ctx, "", "", nil)
	request.Method("PATCH")
	result, err := request.toRequest()

	expected := http.MethodPatch
	if result.Method != expected || err != nil {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodDelete(t *testing.T) {
	request := GetCallBuilder(ctx, "", "", nil)
	request.Method("DELETE")
	result, err := request.toRequest()

	expected := http.MethodDelete
	if result.Method != expected || err != nil {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodEmpty(t *testing.T) {
	request := GetCallBuilder(ctx, "", "", nil)
	request.Method("")
	result, _ := request.toRequest()

	expected := ""
	if result.Method != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestAcceptContentTypeHeaders(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.Accept("acceptHeaderValue")
	request.ContentType("contentTypeHeaderValue")
	result, err := request.toRequest()

	if err != nil || result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" &&
		result.Header.Get(CONTENT_TYPE_HEADER) != "contentTypeHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestHeaders(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.Header(ACCEPT_HEADER, "acceptHeaderValue")
	request.Header("", "empty")
	result, err := request.toRequest()

	if result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" || err != nil {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestCombineHeaders(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.Header(ACCEPT_HEADER, "acceptHeaderValue")
	request.CombineHeaders(map[string]string{CONTENT_TYPE_HEADER: "contentTypeHeaderValue"})
	result, err := request.toRequest()

	if err != nil || result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" &&
		result.Header.Get(CONTENT_TYPE_HEADER) != "contentTypeHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestQueryParam(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.QueryParam("param", "query")
	result, err := request.toRequest()

	if result.URL.RawQuery != "param=query" || err != nil {
		t.Errorf("Failed:\nExpected query param missing")
	}
}

func TestQueryParams(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.QueryParams(map[string]any{"param": "query", "param1": "query"})
	result, err := request.toRequest()

	if !strings.Contains(result.URL.RawQuery, "param=query&param1=query") || err != nil {
		t.Errorf("Failed:\nExpected query params missing")
	}
}

func TestFormData(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	formFields := []FormParam{
		{"param", "form", nil},
		{"param1", "form", nil},
	}
	request.FormData(formFields)
	result, err := request.toRequest()

	if result.Body == nil || err != nil {
		t.Errorf("Failed:\nExpected form data in body")
	}
}

func TestText(t *testing.T) {
	callBuilder := GetCallBuilder(ctx, "GET", "", nil)
	callBuilder.Text("TestString")
	result, err := callBuilder.toRequest()

	stringBuilder := new(strings.Builder)
	io.Copy(stringBuilder, result.Body)

	if !strings.Contains(stringBuilder.String(), "TestString") || err != nil {
		t.Errorf("Failed:\nExpected text in body\n%v", stringBuilder.String())
	}
}

func TestJson(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "", nil)
	request.Json("TestString")
	result, err := request.toRequest()

	stringBuilder := new(strings.Builder)
	_, _ = io.Copy(stringBuilder, result.Body)

	if !strings.Contains(stringBuilder.String(), "TestString") || err != nil {
		t.Errorf("Failed:\nExpected json in body\n%v", stringBuilder.String())
	}
}

func TestFileStream(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "/response/binary", nil)
	request.ContentType("image/png")
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}
	request.FileStream(file)
	_, resp, err := request.CallAsStream()
	if err != nil {
		t.Errorf("Error in CallAsStream: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Failed:\nExpected 200\nGot:%v", resp.StatusCode)
	}
}

func TestFileStreamWithoutHeader(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "/response/binary", nil)
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}
	request.FileStream(file)
	_, resp, err := request.CallAsStream()
	if err != nil {
		t.Errorf("Error in CallAsStream: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Failed:\nExpected 200\nGot:%v", resp.StatusCode)
	}
}

func TestRequestRetryOption(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "/retry", nil)
	request.RequestRetryOption(Disable)

	request2 := GetCallBuilder(ctx, "GET", "/retry", nil)
	if reflect.DeepEqual(request, request2) {
		t.Error("Failed:\nExpected different retry option setting for both requests but got same for both.")
	}
}

func TestContextPropagationInRequests(t *testing.T) {
	key := "Test Key"
	ctx = context.WithValue(ctx, &key, "Test Value")
	request := GetCallBuilder(ctx, "GET", "", nil)
	result, err := request.toRequest()

	if err != nil && result.Context().Value(&key) == "Test Value" {
		t.Errorf("Failed:\nExpected context not found within the request.")
	}
}

func TestRequestCancellation(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	request := GetCallBuilder(ctx, "GET", "", nil)
	_, _, err := request.CallAsJson()

	if err == nil {
		t.Errorf("Failed:\nExpected error due to request cancellation.")
	}
}

func TestError(t *testing.T) {
	codeToErrorMapDefault := map[string]ErrorBuilder[error]{
		"400": {Message: "400 Wrong Input"},
		"5XX": {TemplatedMessage: "Server Error! Server Message: {$response.body#/errorDetail}", Unmarshaller: func(ae ApiError) error { return errors.New(ae.Message) }},
		"0":   {TemplatedMessage: "Error {$statusCode}"},
	}

	codeToErrorMapEmpty := map[string]ErrorBuilder[error]{}

	var tests = []struct {
		name           string
		path           string
		codeToErrorMap map[string]ErrorBuilder[error]
		expected       string
	}{
		{`DefaultErrorMessage`, `/error/400`, codeToErrorMapEmpty, `HTTP Response Not OK`},
		{`StaticErrorMessage`, `/error/400`, codeToErrorMapDefault, `400 Wrong Input`},
		{`DynamicErrorMessageUsingJSONResponseBody`, `/error/500`, codeToErrorMapDefault, `Server Error! Server Message: The server is down at the moment.`},
		{`DynamicErrorMessageWithoutBody`, `/error/404`, codeToErrorMapDefault, `Error 404`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := GetCallBuilder(ctx, "GET", test.path, nil)
			request.AppendErrors(test.codeToErrorMap)

			_, err := request.Call()

			assert.ErrorContains(t, err, test.expected)
		})
	}
}
