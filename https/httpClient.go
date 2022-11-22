package https

import (
	"net/http"
	"time"
)

type HttpClientInterface interface {
	Execute(request *http.Request) *http.Response
}

type HttpClient struct {
	httpClientInstance http.Client
}

func NewHttpClient(opts ...ClientOption) HttpClient {
	httpClient := &HttpClient{
		httpClientInstance: *http.DefaultClient,
	}

	for _, opt := range opts {
		opt(httpClient)
	}

	return *httpClient
}

func (c *HttpClient) Execute(request *http.Request) (*http.Response, error) {
	return c.httpClientInstance.Do(request)
}

type ClientOption func(*HttpClient)

func WithTransport(transport http.RoundTripper) ClientOption {
	return func(h *HttpClient) {
		h.httpClientInstance.Transport = transport
	}
}

func WithTimeout(timeout float64) ClientOption {
	return func(h *HttpClient) {
		h.httpClientInstance.Timeout = time.Duration(timeout * float64(time.Second))
	}
}
