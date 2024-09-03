// Copyright (c) APIMatic. All rights reserved.
package https

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

func (a ApiError) Error() string {
	return fmt.Sprintf("ApiError occured: %v", a.Message)
}

type ErrorBuilder[T error] struct {
	Message          string
	TemplatedMessage string
	Unmarshaller     func(ApiError) T
}

func (eb ErrorBuilder[T]) Build(httpctx HttpContext) error {
	respStatusCode := httpctx.Response.StatusCode
	respHeader := httpctx.Response.Header
	respBodyBytes, _ := httpctx.GetResponseBody()

	message := eb.Message
	if eb.TemplatedMessage != "" {
		message = renderErrorTemplate(eb.TemplatedMessage, respStatusCode, respHeader, respBodyBytes)
	}

	err := ApiError{
		Request:    *httpctx.Request,
		StatusCode: respStatusCode,
		Headers:    respHeader,
		Body:       respBodyBytes,
		Message:    message,
	}

	if eb.Unmarshaller != nil {
		return eb.Unmarshaller(err)
	}

	return err
}

func renderErrorTemplate(tpl string, respStatusCode int, respHeader http.Header, respBytes []byte) string {
	placeholderRegex := `\{\$(.*?)\}`
	re := regexp.MustCompile(placeholderRegex)

	// Extract placeholders into an array of strings
	placeholders := re.FindAllString(tpl, -1)

	renderedVals := []any{}
	for _, placeholder := range placeholders {
		renderedVals = append(renderedVals, renderPlaceholder(placeholder, respStatusCode, respHeader, respBytes))
	}

	// Replace each instance of a placeholder with "%v"
	formattedTpl := re.ReplaceAllString(tpl, "%v")

	return fmt.Sprintf(formattedTpl, renderedVals...)
}

func renderPlaceholder(placeholder string, respStatusCode int, respHeader http.Header, respBytes []byte) any {
	if placeholder == "{$statusCode}" {
		return respStatusCode
	}

	if strings.HasPrefix(placeholder, "{$response.header.") {
		headerName := placeholder[len("{$response.header.") : len(placeholder)-1]
		return respHeader.Get(headerName)
	}

	// Return Response Body as-is
	if placeholder == "{$response.body}" {
		return string(respBytes)
	}

	// Use JSON Pointer to get the desired value from a JSON Response Body
	if strings.HasPrefix(placeholder, "{$response.body#") {
		jsonPtr := placeholder[len("{$response.body#") : len(placeholder)-1]
		if jsonPtr == "" {
			return ""
		}

		return getValueFromJSON(respBytes, jsonPtr)
	}

	return placeholder
}
