package https

import (
	"net/http"
)

type HttpConfigurationOptions func(*HttpConfiguration)

type HttpConfiguration struct {
	timeout            float64
	transport          http.RoundTripper
	retryConfiguration RetryConfiguration
}

func (h *HttpConfiguration) Timeout() float64 {
	return h.timeout
}

func (h *HttpConfiguration) Transport() http.RoundTripper {
	return h.transport
}

func (h *HttpConfiguration) RetryConfiguration() RetryConfiguration {
	return h.retryConfiguration
}

func NewHttpConfiguration(options ...HttpConfigurationOptions) HttpConfiguration {
	httpConfiguration := HttpConfiguration{}

	for _, option := range options {
		option(&httpConfiguration)
	}
	return httpConfiguration
}

func WithTimeout(timeout float64) HttpConfigurationOptions {
	return func(h *HttpConfiguration) {
		h.timeout = timeout
	}
}

func WithTransport(transport http.RoundTripper) HttpConfigurationOptions {
	return func(h *HttpConfiguration) {
		h.transport = transport
	}
}

func WithRetryConfiguration(retryConfig RetryConfiguration) HttpConfigurationOptions {
	return func(h *HttpConfiguration) {
		h.retryConfiguration = retryConfig
	}
}
