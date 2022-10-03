package https

import (
	"net/http"
)

type HttpConfiguration struct {
	Timeout   float64
	Transport http.RoundTripper
}

func DefaultHttpConfiguration() HttpConfiguration {
	return HttpConfiguration{
		Timeout:   0,
		Transport: http.DefaultTransport,
	}
}
