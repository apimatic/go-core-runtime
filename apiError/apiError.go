// Package apiError provides a structure to represent error responses from API calls.
// Copyright (c) APIMatic. All rights reserved.
package apiError

import (
	"fmt"
	"net/http"
)

// ApiError is the base struct for all error responses from the server.
// It holds information about the original HTTP request, the status code, headers, and response body.
type ApiError struct {
	Request    http.Request
	StatusCode int
	Headers    map[string]string
	Body       string
	Message    string
}

// Error implements the Error method for the error interface.
// It returns a string representation of the ApiError instance when used in an error context.
func (a *ApiError) Error() string {
	return fmt.Sprintf("ApiError occured %v", a.Body)
}

type ErrorBuilder[T any] struct {
	Message          string
	TemplatedMessage string
	Builder          func(ApiError) T
}
