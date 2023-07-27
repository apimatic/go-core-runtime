package https

import (
	"net/http"
	"time"
)

// HttpClientInterface defines an interface for an HTTP client
// that can execute HTTP requests and return HTTP responses.
type HttpClientInterface interface {
	Execute(request *http.Request) *http.Response
}

// HttpClient is an implementation of the HttpClientInterface.
type HttpClient struct {
	httpClientInstance http.Client
}

// NewHttpClient creates a new HttpClient with the provided HttpConfiguration and returns it.
// The HttpConfiguration is used to set configurations like timeout, transport and retry configuration settings for the HTTP client.
func NewHttpClient(httpConfig HttpConfiguration) HttpClient {
	client := HttpClient{
		httpClientInstance: *http.DefaultClient,
	}
	client.httpClientInstance.Timeout = time.Duration(httpConfig.Timeout() * float64(time.Second))
	client.httpClientInstance.Transport = httpConfig.Transport()
	return client
}

// Execute sends an HTTP request using the HttpClient and returns the HTTP response or an error.
func (c *HttpClient) Execute(request *http.Request) (*http.Response, error) {
	return c.httpClientInstance.Do(request)
}
