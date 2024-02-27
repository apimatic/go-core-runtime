// Package apiError provides a structure to represent error responses from API calls.
// Copyright (c) APIMatic. All rights reserved.
package https

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-openapi/jsonpointer"
)

// ApiError is the base struct for all error responses from the server.
// It holds information about the original HTTP request, the status code, headers, and response body.
type ApiError struct {
	Request    http.Request
	StatusCode int
	Headers    http.Header
	Body       []byte
	Message    string
}

// Error implements the Error method for the error interface.
// It returns a string representation of the ApiError instance when used in an error context.
func (a *ApiError) Error() string {
	return fmt.Sprintf("ApiError occured: %v", a.Message)
}

type ErrorBuilder[T any] struct {
	Message          string
	TemplatedMessage string
	Unmarshaller     func(ApiError) T
}

func (eb ErrorBuilder[T]) Build(httpctx HttpContext) *ApiError {
	res := httpctx.Response

	message := eb.Message
	if eb.TemplatedMessage != "" {
		message = renderErrorTemplate(eb.TemplatedMessage, *res)
	}

	body, _ := io.ReadAll(res.Body)
	return &ApiError{
		Request:    *httpctx.Request,
		StatusCode: res.StatusCode,
		Headers:    res.Header,
		Body:       body,
		Message:    message,
	}
}

func renderErrorTemplate(tpl string, res http.Response) string {
	placeholderRegex := `\{\$(.*?)\}`
	re := regexp.MustCompile(placeholderRegex)

	// Extract placeholders into an array of strings
	placeholders := re.FindAllString(tpl, -1)

	renderedVals := []any{}
	for _, placeholder := range placeholders {
		renderedVals = append(renderedVals, renderPlaceholder(placeholder, res))
	}

	// Replace each instance of a placeholder with "%v"
	formattedTpl := re.ReplaceAllString(tpl, "%v")

	return fmt.Sprintf(formattedTpl, renderedVals...)
}

func renderPlaceholder(placeholder string, res http.Response) any {
	if placeholder == "{$statusCode}" {
		return res.StatusCode
	}

	if strings.HasPrefix(placeholder, "{$response.header.") {
		headerName := placeholder[len("{$response.header.") : len(placeholder)-1]
		return res.Header.Get(headerName)
	}

	if placeholder == "{$response.body}" {
		serializedBody, err := io.ReadAll(res.Body)

		if err != nil {
			return ""
		}
		return string(serializedBody)
	}

	// Use JSON Pointer to get the desired value
	if strings.HasPrefix(placeholder, "{$response.body#") {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ""
		}

		var jsonBody any
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			return ""
		}

		jsonPtr := placeholder[len("{$response.body#") : len(placeholder)-1]
		if jsonPtr == "" {
			return ""
		}

		p, err := jsonpointer.New(jsonPtr)
		if err != nil {
			return ""
		}

		val, kind, err := p.Get(jsonBody)
		if err != nil {
			return ""
		}

		switch kind {
		case reflect.Map:
			obj, err := json.Marshal(val)

			if err != nil {
				return ""
			}
			return string(obj)
		}

		return val
	}

	return placeholder
}
