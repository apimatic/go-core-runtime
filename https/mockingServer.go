package https

import (
	"net/http"
	"net/http/httptest"
)

// GetTestingServer creates and returns an httptest.Server instance for testing purposes.
func GetTestingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Path == "/response/integer" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`4`))

			} else if r.URL.Path == "/template/abc/def" || r.URL.Path == "/template/1/2/3/4/5" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"passed": true,
				"message": "It's a hit!",`))
			} else if r.URL.Path == "/response/binary" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"passed": true,
				"message": "It's a hit!",`))
			}
		} else if r.Method == "POST" {
			if r.URL.Path == "/form/string" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`4`))
			}
		}
	}))
}
