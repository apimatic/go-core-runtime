package https

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMethod(t *testing.T) {
	expected := "ApiError occured Server Error"
	result := ApiError{
		StatusCode: 500,
		Body:       "Server Error",
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
	tpl := "Error: Date {$response.headers.Date}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Date Thu, 22 Feb 2024 06:01:57 GMT", actual)
}

func TestCorrectMessageWhenDynamicErrorMessageWithResponseBody(t *testing.T) {
	body := []byte(`{
		"Error": [
		  {
			"Code": 1,
			"Type": "Critical"
		  }
		]
	  }`)
	res := http.Response{
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	tpl := "Error: Code {$response.body#/Error/0/Code}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Code 1", actual)
}

func TestCorrectMessageWhenDynamicErrorMessageWithResponseBodyLargeValue(t *testing.T) {
	body := []byte(`{
		"Error": [
		  {
			"Code": 100000000000000000,
			"Type": "Critical"
		  }
		]
	  }`)
	res := http.Response{
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	tpl := "Error: Code {$response.body#/Error/0/Code}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Code 100000000000000000", actual)
}

func TestEmptyStringWhenDynamicErrorMessageWithMissingResponseHeader(t *testing.T) {
	h := http.Header{}
	res := http.Response{
		Header: h,
	}
	tpl := "Error: Date {$response.headers.Date}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Date ", actual)
}

func TestEmptyStringWhenDynamicErrorMessageWithResponseBodyPropertyMissing(t *testing.T) {
	body := []byte(`{
		"Error": [
		  {
			"Type": "Critical"
		  }
		]
	  }`)
	res := http.Response{
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	tpl := "Error: Code {$response.body#/Error/0/Code}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Code ", actual)
}

func TestEmptyStringWhenDynamicErrorMessageWithInvalidJsonPointer(t *testing.T) {
	body := []byte(`{
		"Error": [
		  {
			"Code": 1,
			"Type": "Critical"
		  }
		]
	  }`)
	res := http.Response{
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	tpl := "Error: Code {$response.body##\\#}"

	actual := renderErrorTemplate(tpl, res)

	assert.Equal(t, "Error: Code ", actual)
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
