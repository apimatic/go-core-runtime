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
	Request    http.Request      `json:"Request"`
	StatusCode int               `json:"StatusCode"`
	Headers    map[string]string `json:"Headers"`
	Body       string            `json:"Body"`
	Message    string
}

// NewApiError is the constructor function for ApiError.
// It creates and returns a pointer to an ApiError instance with the given status code and response body.
func NewApiError(
	errorType string,
	message string) *ApiError {
	return &ApiError{
		Message: message,
	}
}

func NewApiErrorTemplated(
	errorType string,
	message string) *ApiError {
	return &ApiError{
		TemplatedMessage: message,
	}
}

// Error implements the Error method for the error interface.
// It returns a string representation of the ApiError instance when used in an error context.
func (a *ApiError) Error() string {
	return fmt.Sprintf("ApiError occured %v", a.Body)
}
