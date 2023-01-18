package https

import "net/http"

type ApiResponse[T any] struct {
	Data     T              `json:"data"`
	Response *http.Response `json:"response"`
}

func NewApiResponse[T any](data T, response *http.Response) ApiResponse[T] {
	apiResponse := ApiResponse[T]{
		Data:     data,
		Response: response,
	}
	return apiResponse
}
