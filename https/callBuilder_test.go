package https

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var ctx context.Context = context.Background()

func GetCallBuilder(ctx context.Context, method, path string, auth Authenticator) CallBuilder {
	return GetCallBuilder1(ctx, method, path, map[string]Authenticator{ "default" : auth})
}
func GetCallBuilder1(ctx context.Context, method, path string, auths map[string]Authenticator) CallBuilder {
	
	client := NewHttpClient(NewHttpConfiguration())
	callBuilder := CreateCallBuilderFactory(
		func(server string) string {
			return GetTestingServer().URL
		},
		auths,
		client,
		NewRetryConfiguration(),
	)

	return callBuilder(ctx, method, path)
}

func RequestAuthentication() Authenticator {
	return func(requiresAuth bool) HttpInterceptor {
		return PassThroughInterceptor
	}
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
	request := GetCallBuilder(ctx, "GET", "/template/%s", nil)
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
	request := GetCallBuilder(ctx, "GET", "/template/%s", nil)
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
	request.QueryParams(map[string]interface{}{"param": "query", "param1": "query"})
	result, err := request.toRequest()

	if !strings.Contains(result.URL.RawQuery, "param=query&param1=query") || err != nil {
		t.Errorf("Failed:\nExpected query params missing")
	}
}

func TestAuthenticate(t *testing.T) {
	request := GetCallBuilder(ctx, "GET", "/auth", RequestAuthentication())
	request.Authenticate([]map[string]bool{{ "oAuthBearerToken": true }})
}

func TestAuthenticates(t *testing.T) {
	auths := map[string]Authenticator{
		"oAuthBearerToken" : RequestAuthentication(),
		"basicAuth" : RequestAuthentication(),
		"apiKey" : RequestAuthentication(),
		"apiHeader" : RequestAuthentication(),
	}
	request := GetCallBuilder1(ctx, "GET", "/auth", auths)
	request.Authenticate([]map[string]bool {
		{ "oAuthBearerToken": false }, { "basicAuth": true, "apiKey": true, "apiHeader": true },
	})
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
	io.Copy(stringBuilder, result.Body)

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
