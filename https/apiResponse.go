package https

import "net/http"

type ApiResponse [T any] struct {

	Data T 					`json:"data"`
	Response *http.Response `json:"response"`
	Error  error 			`json:"error"`
}

func NewApiResponse[T any] (data T, response *http.Response, err error) ApiResponse[T] {
	
	apiResponse := ApiResponse[T]{
		Data: data,
		Response: response,
		Error: err,
	}
	return apiResponse
}