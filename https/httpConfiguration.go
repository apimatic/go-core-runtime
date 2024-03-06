package https

import (
	"net/http"
)

// HttpConfigurationOptions is a function type that takes a pointer to HttpConfiguration as input and modifies it.
type HttpConfigurationOptions func(*HttpConfiguration)

// HttpConfiguration holds the configuration options for the HTTP client.
type HttpConfiguration struct {
	timeout            float64
	transport          http.RoundTripper
	retryConfiguration RetryConfiguration
}

// Timeout returns the configured timeout value for the HTTP client.
func (h *HttpConfiguration) Timeout() float64 {
	return h.timeout
}

// Transport returns the configured custom HTTP transport for the HTTP client.
func (h *HttpConfiguration) Transport() http.RoundTripper {
	return h.transport
}

// RetryConfiguration returns the configured retry configuration for the HTTP client.
func (h *HttpConfiguration) RetryConfiguration() RetryConfiguration {
	return h.retryConfiguration
}

// NewHttpConfiguration creates a new HttpConfiguration with the provided options and returns it.
// The options parameter allows setting different configuration options for the HTTP client.
func NewHttpConfiguration(options ...HttpConfigurationOptions) HttpConfiguration {
	httpConfiguration := HttpConfiguration{}

	for _, option := range options {
		option(&httpConfiguration)
	}
	return httpConfiguration
}

// WithTimeout sets the timeout for the HTTP client.
func WithTimeout(timeout float64) HttpConfigurationOptions {
	return func(h *HttpConfiguration) {
		h.timeout = timeout
	}
}

// WithTransport sets the custom HTTP transport for the HTTP client.
func WithTransport(transport http.RoundTripper) HttpConfigurationOptions {
	return func(h *HttpConfiguration) {
		h.transport = transport
	}
}

// WithRetryConfiguration sets the retry configuration for the HTTP client.
func WithRetryConfiguration(retryConfig RetryConfiguration) HttpConfigurationOptions {
	return func(h *HttpConfiguration) {
		h.retryConfiguration = retryConfig
	}
}
