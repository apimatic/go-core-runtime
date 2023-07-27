// Package https provides utilities and structures for handling
// API requests and responses using the HTTP protocol.
// Copyright (c) APIMatic. All rights reserved.
package https

import "net/http"

// ApiResponse is a generic struct that represents an API response containing data and the HTTP response.
// The `Data` field holds the data of any type `T` returned by the API.
// The `Response` field contains the underlying HTTP response associated with the API call.
type ApiResponse[T any] struct {
	Data     T              `json:"data"`
	Response *http.Response `json:"response"`
}

// NewApiResponse creates a new instance of ApiResponse.
// It takes the `data` of type `T` and the `response` as parameters, and returns an ApiResponse[T] struct.
func NewApiResponse[T any](data T, response *http.Response) ApiResponse[T] {
	apiResponse := ApiResponse[T]{
		Data:     data,
		Response: response,
	}
	return apiResponse
}
