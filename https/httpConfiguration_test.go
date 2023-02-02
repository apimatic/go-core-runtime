package https

import (
	"net/http"
	"reflect"
	"testing"
)

func getHttpConfig() HttpConfiguration {
	return NewHttpConfiguration(
		WithTimeout(100),
		WithTransport(http.DefaultTransport),
		WithRetryConfiguration(getRetryConfiguration()))
}
func TestDefaultConfig(t *testing.T) {
	got := getHttpConfig()

	if got.Timeout() != 100 {
		t.Errorf("Failed:\nExpected Timeout: %v\nGot: %v", 100, got.Timeout())
	}
	if !reflect.DeepEqual(got.Transport(), http.DefaultTransport) {
		t.Errorf("Failed:\nExpected Transport: %v\nGot: %v", http.DefaultTransport, got.Transport())
	}
	if !reflect.DeepEqual(got.RetryConfiguration(), getRetryConfiguration()) {
		t.Errorf("Failed:\nExpected Retry Configuration: %v\nGot: %v", getRetryConfiguration(), got.RetryConfiguration())
	}
}
