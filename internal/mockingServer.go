package internal

import (
	"net/http"
	"net/http/httptest"
)

// GetTestingServer creates and returns a httptest.Server instance for testing purposes.
func GetTestingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			switch r.URL.Path {
			case "/response/integer":
				_, _ = w.Write([]byte(`4`))
			case "/template/abc/def", "/template/1/2/3/4/5", "/response/binary":
				_, _ = w.Write([]byte(`"passed": true,
				"message": "It's a hit!",`))
			case "/error/400":
				w.WriteHeader(http.StatusBadRequest)
			case "/error/500":
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"errorDetail":"The server is down at the moment."}`))
			case "/error/404":
				w.WriteHeader(http.StatusNotFound)
			}
		case "POST":
			switch r.URL.Path {
			case "/form/string":
				_, _ = w.Write([]byte(`4`))
			}
		}
	}))
}
