package https

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"time"
)

type RequestRetryOption int

const (
	Default RequestRetryOption = iota
	Enable
	Disable
)

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

type RetryConfigurationOptions func(*RetryConfiguration)

type RetryConfiguration struct {
	maxRetryAttempts       int64
	retryOnTimeout         bool
	retryInterval          time.Duration
	maximumRetryWaitTime   time.Duration
	backoffFactor          int64
	httpStatusCodesToRetry []int64
	httpMethodsToRetry     []string
}

func NewRetryConfiguration(options ...RetryConfigurationOptions) RetryConfiguration {
	retryConfig := RetryConfiguration{}

	for _, option := range options {
		option(&retryConfig)
	}
	return retryConfig
}

func WithMaxRetryAttempts(maxRetryAttempts int64) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.maxRetryAttempts = maxRetryAttempts
	}
}

func WithRetryOnTimeout(retryOnTimeout bool) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.retryOnTimeout = retryOnTimeout
	}
}

func WithRetryInterval(retryInterval time.Duration) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.retryInterval = retryInterval
	}
}

func WithMaximumRetryWaitTime(maximumRetryWaitTime time.Duration) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.maximumRetryWaitTime = maximumRetryWaitTime
	}
}

func WithBackoffFactor(backoffFactor int64) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.backoffFactor = backoffFactor
	}
}

func WithHttpStatusCodesToRetry(httpStatusCodesToRetry []int64) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.httpStatusCodesToRetry = httpStatusCodesToRetry
	}
}

func WithHttpMethodsToRetry(httpMethodsToRetry []string) RetryConfigurationOptions {
	return func(r *RetryConfiguration) {
		r.httpMethodsToRetry = httpMethodsToRetry
	}
}

func (r *RetryConfiguration) MaxRetryAttempts() int64 {
	return r.maxRetryAttempts
}

func (r *RetryConfiguration) MaximumRetryWaitTime() time.Duration {
	return r.maximumRetryWaitTime
}

func (r *RetryConfiguration) BackoffFactor() int64 {
	return r.backoffFactor
}

func (r *RetryConfiguration) RetryInterval() time.Duration {
	return r.retryInterval
}

func (r *RetryConfiguration) RetryOnTimeout() bool {
	return r.retryOnTimeout
}

func (r *RetryConfiguration) HttpMethodsToRetry() []string {
	return r.httpMethodsToRetry
}

func (r *RetryConfiguration) HttpStatusCodesToRetry() []int64 {
	return r.httpStatusCodesToRetry
}

func (rc *RetryConfiguration) containsHttpStatusCodesToRetry(httpStatusCode int64) bool {
	for _, val := range rc.httpStatusCodesToRetry {
		if val == httpStatusCode {
			return true
		}
	}
	return false
}

func (rc *RetryConfiguration) containsHttpMethodsToRetry(httpMethod string) bool {
	for _, v := range rc.httpMethodsToRetry {
		if v == httpMethod {
			return true
		}
	}
	return false
}

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

func getRetryAfterInSeconds(headers http.Header) time.Duration {
	if s, ok := headers[http.CanonicalHeaderKey("retry-after")]; ok {
		if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
			return time.Duration(sleep)
		}
	}
	return 0
}

func defaultBackoff(retryInterval, maxWaitTime, retryAfter time.Duration, backoffFactor, retryCount int64) time.Duration {
	waitTime := (math.Pow(float64(backoffFactor), float64(retryCount)) * float64(retryInterval))
	sleep := math.Max(waitTime, float64(retryAfter))
	if time.Duration(sleep) <= maxWaitTime {
		return time.Duration(sleep)
	}
	return 0
}
