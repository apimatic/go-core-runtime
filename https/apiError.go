// Copyright (c) APIMatic. All rights reserved.
package https

import (
	"fmt"
	"io"
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
	res := httpctx.Response

	message := eb.Message
	if eb.TemplatedMessage != "" {
		message = renderErrorTemplate(eb.TemplatedMessage, *res)
	}

	body, _ := io.ReadAll(res.Body)
	err := ApiError{
		Request:    *httpctx.Request,
		StatusCode: res.StatusCode,
		Headers:    res.Header,
		Body:       body,
		Message:    message,
	}

	if eb.Unmarshaller != nil {
		return eb.Unmarshaller(err)
	}

	return err
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

	// Return Response Body as-is
	if placeholder == "{$response.body}" {
		serializedBody, err := io.ReadAll(res.Body)
		if err != nil {
			return ""
		}

		return string(serializedBody)
	}

	// Use JSON Pointer to get the desired value from a JSON Response Body
	if strings.HasPrefix(placeholder, "{$response.body#") {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ""
		}

		jsonPtr := placeholder[len("{$response.body#") : len(placeholder)-1]
		if jsonPtr == "" {
			return ""
		}

		return getValueFromJSON(body, jsonPtr)
	}

	return placeholder
}
