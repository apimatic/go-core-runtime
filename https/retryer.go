package https

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"time"
)

// RequestRetryOption represents the type for request retry options.
type RequestRetryOption int

// Constants for different request retry options.
const (
	Default RequestRetryOption = iota
	Enable
	Disable
)

// String returns the string representation of the RequestRetryOption.
func (r RequestRetryOption) String() string {
	switch r {
	case Enable:
		return "enable"
	case Disable:
		return "disable"
	default:
		return "default"
	}
}

// RetryConfigurationOptions represents a function that modifies RetryConfiguration settings.
type RetryConfigurationOptions func(*RetryConfiguration)

// RetryConfiguration contains the settings for request retry behavior.
type RetryConfiguration struct {
	maxRetryAttempts       int64
	retryOnTimeout         bool
	retryInterval          time.Duration
	maximumRetryWaitTime   time.Duration
	backoffFactor          int64
	httpStatusCodesToRetry []int64
	httpMethodsToRetry     []string
}

// NewRetryConfiguration creates a new RetryConfiguration instance with the provided options.
func NewRetryConfiguration(options ...RetryConfigurationOptions) RetryConfiguration {
	retryConfig := RetryConfiguration{}

	for _, option := range options {
		option(&retryConfig)
	}
	return retryConfig
}

// WithMaxRetryAttempts sets the maximum number of retry attempts allowed.
func WithMaxRetryAttempts(maxRetryAttempts int64) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.maxRetryAttempts = maxRetryAttempts
	}
}

// WithRetryOnTimeout sets whether to retry on timeouts.
func WithRetryOnTimeout(retryOnTimeout bool) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.retryOnTimeout = retryOnTimeout
	}
}

// WithRetryInterval sets the interval between retries.
func WithRetryInterval(retryInterval time.Duration) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.retryInterval = retryInterval
	}
}

// WithMaximumRetryWaitTime sets the maximum wait time before giving up retrying.
func WithMaximumRetryWaitTime(maximumRetryWaitTime time.Duration) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.maximumRetryWaitTime = maximumRetryWaitTime
	}
}

// WithBackoffFactor sets the backoff factor for exponential backoff.
func WithBackoffFactor(backoffFactor int64) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.backoffFactor = backoffFactor
	}
}

// WithHttpStatusCodesToRetry sets the list of HTTP status codes to retry on.
func WithHttpStatusCodesToRetry(httpStatusCodesToRetry []int64) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.httpStatusCodesToRetry = httpStatusCodesToRetry
	}
}

// WithHttpMethodsToRetry sets the list of HTTP methods to retry on.
func WithHttpMethodsToRetry(httpMethodsToRetry []string) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.httpMethodsToRetry = httpMethodsToRetry
	}
}

// MaxRetryAttempts returns the maximum number of retry attempts allowed.
func (r *RetryConfiguration) MaxRetryAttempts() int64 {
	return r.maxRetryAttempts
}

// MaximumRetryWaitTime returns the maximum wait time before giving up retrying.
func (r *RetryConfiguration) MaximumRetryWaitTime() time.Duration {
	return r.maximumRetryWaitTime
}

// BackoffFactor returns the backoff factor for exponential backoff.
func (r *RetryConfiguration) BackoffFactor() int64 {
	return r.backoffFactor
}

// RetryInterval returns the interval between retries.
func (r *RetryConfiguration) RetryInterval() time.Duration {
	return r.retryInterval
}

// RetryOnTimeout returns whether to retry on timeouts.
func (r *RetryConfiguration) RetryOnTimeout() bool {
	return r.retryOnTimeout
}

// HttpMethodsToRetry returns the list of HTTP methods to retry on.
func (r *RetryConfiguration) HttpMethodsToRetry() []string {
	return r.httpMethodsToRetry
}

// HttpStatusCodesToRetry returns the list of HTTP status codes to retry on.
func (r *RetryConfiguration) HttpStatusCodesToRetry() []int64 {
	return r.httpStatusCodesToRetry
}

// containsHttpStatusCodesToRetry checks if the given HTTP status code exists in the list of status codes to retry on.
func (rc *RetryConfiguration) containsHttpStatusCodesToRetry(httpStatusCode int64) bool {
	for _, val := range rc.httpStatusCodesToRetry {
		if val == httpStatusCode {
			return true
		}
	}
	return false
}

// containsHttpMethodsToRetry checks if the given HTTP method exists in the list of methods to retry on.
func (rc *RetryConfiguration) containsHttpMethodsToRetry(httpMethod string) bool {
	for _, v := range rc.httpMethodsToRetry {
		if v == httpMethod {
			return true
		}
	}
	return false
}

// GetRetryWaitTime calculates the wait time for the next retry attempt.
func (rc *RetryConfiguration) GetRetryWaitTime(
	maxWaitTime time.Duration,
	retryCount int64,
	response *http.Response,
	timeoutError error) time.Duration {
	retry := false
	var retryAfter, retryWaitTime time.Duration

	if retryCount < rc.maxRetryAttempts {
		if err, ok := timeoutError.(net.Error); ok && err.Timeout() {
			retry = rc.retryOnTimeout
		} else if response != nil && response.Header != nil {
			retryAfter = getRetryAfterInSeconds(response.Header)
			if retryAfter > 0 || rc.containsHttpStatusCodesToRetry(int64(response.StatusCode)) {
				retry = true
			}
		}

		if retry {
			retryWaitTime = defaultBackoff(rc.retryInterval, maxWaitTime, retryAfter, rc.backoffFactor, retryCount)
		}
	}
	return retryWaitTime
}

// ShouldRetry determines if the request should be retried based on the RetryConfiguration and request HTTP method.
func (rc *RetryConfiguration) ShouldRetry(retryRequestOption RequestRetryOption, httpMethod string) bool {
	switch retryRequestOption.String() {
	default:
		if rc.maxRetryAttempts > 0 && httpMethod != "" && rc.containsHttpMethodsToRetry(httpMethod) {
			return true
		}
		return false
	case "enable":
		return true
	case "disable":
		return false
	}
}

// getRetryAfterInSeconds extracts and returns the Retry-After duration from the response headers in seconds.
func getRetryAfterInSeconds(headers http.Header) time.Duration {
	if s, ok := headers[http.CanonicalHeaderKey("retry-after")]; ok {
		if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
			return time.Duration(sleep)
		}
	}
	return 0
}

// defaultBackoff calculates the backoff time for exponential backoff based on the retry configuration.
func defaultBackoff(retryInterval, maxWaitTime, retryAfter time.Duration, backoffFactor, retryCount int64) time.Duration {
	waitTime := (math.Pow(float64(backoffFactor), float64(retryCount)) * float64(retryInterval))
	sleep := math.Max(waitTime, float64(retryAfter))
	if time.Duration(sleep) <= maxWaitTime {
		return time.Duration(sleep)
	}
	return 0
}
