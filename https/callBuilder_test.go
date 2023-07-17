package https

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func GetCallBuilder(method, path string, auth Authenticator) CallBuilder {
	client := NewHttpClient(NewHttpConfiguration())
	ctx := context.WithValue(context.Background(), "Test Key", "Test Value")
	callBuilder := CreateCallBuilderFactory(
		func(server string) string {
			return GetTestingServer().URL
		},
		auth,
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
	request := GetCallBuilder("GET", "//response/", nil)
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

func TestAppendPathEmptyPath(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
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
	request := GetCallBuilder("GET", "/template/%s", nil)
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
	request := GetCallBuilder("GET", "/template/%s", nil)
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
	request := GetCallBuilder("", "/response/integer", nil)
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
	request := GetCallBuilder("", "/response/integer", nil)
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
	request := GetCallBuilder("", "/form/string", nil)
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
	request := GetCallBuilder("", "", nil)
	request.Method("PUT")
	result, err := request.toRequest()

	expected := http.MethodPut
	if result.Method != expected || err != nil {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodPatch(t *testing.T) {
	request := GetCallBuilder("", "", nil)
	request.Method("PATCH")
	result, err := request.toRequest()

	expected := http.MethodPatch
	if result.Method != expected || err != nil {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodDelete(t *testing.T) {
	request := GetCallBuilder("", "", nil)
	request.Method("DELETE")
	result, err := request.toRequest()

	expected := http.MethodDelete
	if result.Method != expected || err != nil {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodEmpty(t *testing.T) {
	request := GetCallBuilder("", "", nil)
	request.Method("")
	result, _ := request.toRequest()

	expected := ""
	if result.Method != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestAcceptContentTypeHeaders(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Accept("acceptHeaderValue")
	request.ContentType("contentTypeHeaderValue")
	result, err := request.toRequest()

	if err != nil || result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" &&
		result.Header.Get(CONTENT_TYPE_HEADER) != "contentTypeHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestHeaders(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Header(ACCEPT_HEADER, "acceptHeaderValue")
	request.Header("", "empty")
	result, err := request.toRequest()

	if result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" || err != nil {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestCombineHeaders(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Header(ACCEPT_HEADER, "acceptHeaderValue")
	request.CombineHeaders(map[string]string{CONTENT_TYPE_HEADER: "contentTypeHeaderValue"})
	result, err := request.toRequest()

	if err != nil || result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" &&
		result.Header.Get(CONTENT_TYPE_HEADER) != "contentTypeHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestQueryParam(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.QueryParam("param", "query")
	result, err := request.toRequest()

	if result.URL.RawQuery != "param=query" || err != nil {
		t.Errorf("Failed:\nExpected query param missing")
	}
}

func TestQueryParams(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.QueryParams(map[string]interface{}{"param": "query", "param1": "query"})
	result, err := request.toRequest()

	if !strings.Contains(result.URL.RawQuery, "param=query&param1=query") || err != nil {
		t.Errorf("Failed:\nExpected query params missing")
	}
}

func TestAuthenticate(t *testing.T) {
	request := GetCallBuilder("GET", "/auth", RequestAuthentication())
	request.Authenticate(true)
}

func TestFormData(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
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
	callBuilder := GetCallBuilder("GET", "", nil)
	callBuilder.Text("TestString")
	result, err := callBuilder.toRequest()

	stringBuilder := new(strings.Builder)
	io.Copy(stringBuilder, result.Body)

	if !strings.Contains(stringBuilder.String(), "TestString") || err != nil {
		t.Errorf("Failed:\nExpected text in body\n%v", stringBuilder.String())
	}
}

func TestJson(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Json("TestString")
	result, err := request.toRequest()

	stringBuilder := new(strings.Builder)
	io.Copy(stringBuilder, result.Body)

	if !strings.Contains(stringBuilder.String(), "TestString") || err != nil {
		t.Errorf("Failed:\nExpected json in body\n%v", stringBuilder.String())
	}
}

func TestFileStream(t *testing.T) {
	request := GetCallBuilder("GET", "/response/binary", nil)
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
	request := GetCallBuilder("GET", "/response/binary", nil)
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
	request := GetCallBuilder("GET", "/retry", nil)
	request.RequestRetryOption(Disable)

	request2 := GetCallBuilder("GET", "/retry", nil)
	if reflect.DeepEqual(request, request2) {
		t.Error("Failed:\nExpected different retry option setting for both requests but got same for both.")
	}
}

func TestContextPropagationInRequests(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	result, err := request.toRequest()

	if err != nil && result.Context().Value("Test Key") == "Test Value" {
		t.Errorf("Failed:\nExpected context not found within the request.")
	}
}
