package apiError

import (
	"fmt"
	"net/http"
)

// This is the base struct for all exceptions that represent an error response from the server.
type ApiError struct {
	Request    http.Request      `json:"Request"`
	StatusCode int               `json:"StatusCode"`
	Headers    map[string]string `json:"Headers"`
	Body       string            `json:"Body"`
}

// Constructor for ApiError.
func NewApiError(
	statusCode int,
	body string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Body:       body,
	}
}

// Implementing the Error method for the error interface.
func (a *ApiError) Error() string {
	return fmt.Sprintf("ApiError occured %v", a.Body)
}
