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

func NewHttpClient(httpConfig HttpConfiguration) HttpClient {
	client := HttpClient{
		httpClientInstance: *http.DefaultClient,
	}
	client.httpClientInstance.Timeout = time.Duration(httpConfig.Timeout() * float64(time.Second))
	client.httpClientInstance.Transport = httpConfig.Transport()
	return client
}

func (c *HttpClient) Execute(request *http.Request) (*http.Response, error) {
	return c.httpClientInstance.Do(request)
}
