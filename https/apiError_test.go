package https

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/apimatic/go-core-runtime/assert"
)

const mockJSONResponseBody = `{
	"Error": [
	  {
		"Code": 1,
		"Type": "Critical"
	  }
	]
  }`

func getMockResponseWithJSONBody(jsonStr string) http.Response {
	body := []byte(jsonStr)

	buffer := new(bytes.Buffer)

	if err := json.Compact(buffer, body); err != nil {
		log.Fatal(err)
	}
	return http.Response{
		Body: io.NopCloser(bytes.NewReader(buffer.Bytes())),
	}
}

func TestErrorMethod(t *testing.T) {
	expected := "ApiError occured: Server Error"
	result := ApiError{
		StatusCode: 500,
		Message:    "Server Error",
	}
	if result.Error() != expected {
		t.Errorf("Failed:\nExpected: %v\nGot: %v", expected, result.Error())
	}
}

func TestCorrectMessageWhenDynamicErrorMessageWithStatusCode(t *testing.T) {
	res := http.Response{
		StatusCode: 500,
	}
	tpl := "Error: Status Code {$statusCode}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Status Code 500", actual)
}

func TestCorrectMessageWhenDynamicErrorMessageWithResponseHeader(t *testing.T) {
	h := http.Header{}
	h.Set("Date", "Thu, 22 Feb 2024 06:01:57 GMT")
	res := http.Response{
		Header: h,
	}
	tpl := "Error: Date {$response.header.Date}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Date Thu, 22 Feb 2024 06:01:57 GMT", actual)
}

func TestDynamicErrorMessageWithResponseBody(t *testing.T) {
	var tests = []struct {
		name     string
		tpl      string
		expected string
	}{
		{`FormattedMessageWhenIntValue`, `Error: Code {$response.body#/Error/0/Code}`, `Error: Code 1`},
		{`JSONStringWhenObjectValue`, `Error: {$response.body#/Error/0}`, `Error: {"Code":1,"Type":"Critical"}`},
		{`RawResponseBodyWhenNoJSONPointer`, `Error: {$response.body}`, `Error: {"Error":[{"Code":1,"Type":"Critical"}]}`},
		{`EmptyStringWhenInvalidJSONPointer`, `Error: Code {$response.body##\\#}`, `Error: Code `},
		{`EmptyStringWhenEmptyJSONPointer`, `Error: Code {$response.body#}`, `Error: Code `},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := getMockResponseWithJSONBody(mockJSONResponseBody)
			actual := renderErrorTemplate(test.tpl, res)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestEmptyStringWhenDynamicErrorMessageWithMissingResponseHeader(t *testing.T) {
	h := http.Header{}
	res := http.Response{
		Header: h,
	}
	tpl := "Error: Date {$response.header.Date}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Date ", actual)
}

func TestEmptyStringWhenDynamicErrorMessageWithResponseBodyPropertyMissing(t *testing.T) {
	res := getMockResponseWithJSONBody(`{
		"Error": [
		  {
			"Type": "Critical"
		  }
		]
	  }`)
	tpl := "Error: Code {$response.body#/Error/0/Code}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Code ", actual)
}

func TestEmptyStringWhenDynamicErrorMessageWithInvalidJSONInResponseBody(t *testing.T) {
	res := http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(`"invalidJson"}`))),
	}
	tpl := "Error: {$response.body#/Id}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: ", actual)
}

func TestReturnWithoutChangesWhenDynamicErrorMessageWithInvalidPlaceholder(t *testing.T) {
	res := http.Response{}
	tpl := "Error: Code {$something.something}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, tpl, actual)
}

func TestReturnWithoutChangesWhenDynamicErrorMessageWithNoTemplates(t *testing.T) {
	res := http.Response{}
	tpl := "An error occurred."

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, tpl, actual)
}
