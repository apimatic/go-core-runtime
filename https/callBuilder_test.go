package https

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func GetCallBuilder(method, path string, auth Authenticator) CallBuilder {
	client := NewHttpClient()
	callBuilder := CreateCallBuilderFactory(
		func(server string) string { return GetTestingServer().URL }, auth, client)

	return callBuilder(method, path)
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

	expected := 200

	if response.StatusCode != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, response)
	}
}

func TestMethodPut(t *testing.T) {
	request := GetCallBuilder("", "", nil)
	request.Method("PUT")
	result, _ := request.toRequest()

	expected := http.MethodPut
	if result.Method != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodPatch(t *testing.T) {
	request := GetCallBuilder("", "", nil)
	request.Method("PATCH")
	result, _ := request.toRequest()

	expected := http.MethodPatch
	if result.Method != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

func TestMethodDelete(t *testing.T) {
	request := GetCallBuilder("", "", nil)
	request.Method("DELETE")
	result, _ := request.toRequest()

	expected := http.MethodDelete
	if result.Method != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Method)
	}
}

// func TestMethodEmpty(t *testing.T) {
// 	request := GetCallBuilder("", "", nil)
// 	err := request.Method("")
// 	if err == nil {
// 				t.Errorf("The code should get error because Invalid HTTP method given!")

// 	}
// }

func TestAcceptContentTypeHeaders(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Accept("acceptHeaderValue")
	request.ContentType("contentTypeHeaderValue")
	result, _ := request.toRequest()

	if result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" &&
		result.Header.Get(CONTENT_TYPE_HEADER) != "contentTypeHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestHeaders(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Header(ACCEPT_HEADER, "acceptHeaderValue")
	request.Header("", "empty")
	result, _ := request.toRequest()

	if result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestCombineHeaders(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Header(ACCEPT_HEADER, "acceptHeaderValue")
	request.CombineHeaders(map[string]string{CONTENT_TYPE_HEADER: "contentTypeHeaderValue"})
	result, _ := request.toRequest()

	if result.Header.Get(ACCEPT_HEADER) != "acceptHeaderValue" &&
		result.Header.Get(CONTENT_TYPE_HEADER) != "contentTypeHeaderValue" {
		t.Errorf("Failed:\nExpected headers not received")
	}
}

func TestQueryParam(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.QueryParam("param", "query")
	result, _ := request.toRequest()

	if result.URL.RawQuery != "param=query" {
		t.Errorf("Failed:\nExpected query param missing")
	}
}

func TestQueryParams(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.QueryParams(map[string]interface{}{"param": "query", "param1": "query"})
	result, _ := request.toRequest()

	if result.URL.RawQuery != "param=query&param1=query" {
		t.Errorf("Failed:\nExpected query params missing")
	}
}

func TestAuthenticate(t *testing.T) {
	request := GetCallBuilder("GET", "/auth", RequestAuthentication())
	request.Authenticate(true)
}

func TestFormData(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.FormData(map[string]interface{}{"param": "form", "param1": "form"})
	result, _ := request.toRequest()

	if result.Body == nil {
		t.Errorf("Failed:\nExpected form data in body")
	}
}

func TestText(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Text("Body Text")
	result, _ := request.toRequest()

	stringBuilder := new(strings.Builder)
	io.Copy(stringBuilder, result.Body)

	if stringBuilder.String() != "Body Text" {
		t.Errorf("Failed:\nExpected text in body")
	}
}

func TestJson(t *testing.T) {
	request := GetCallBuilder("GET", "", nil)
	request.Json("Json")
	result, _ := request.toRequest()

	stringBuilder := new(strings.Builder)
	io.Copy(stringBuilder, result.Body)

	if stringBuilder.String() != `"Json"` {
		t.Errorf("Failed:\nExpected json in body")
	}
}

// func TestJsonPanic(t *testing.T) {
// 	request := GetCallBuilder("", "", nil)
// 	err := request.Json(math.Inf(2))
// 	if err == nil {
// 		t.Errorf("The code should get error because request is empty.")
// 	}
// }

func TestFileStream(t *testing.T) {
	request := GetCallBuilder("GET", "/response/binary", nil)
	request.ContentType("image/png")
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		err = fmt.Errorf("GetFile failed: %v", err)
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
		err = fmt.Errorf("GetFile failed: %v", err)
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
