package https

import (
	"net/http"
	"net/http/httptest"
)

// GetTestingServer creates and returns an httptest.Server instance for testing purposes.
func GetTestingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			switch r.URL.Path {
			case "/response/integer":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`4`))
			case "/template/abc/def", "/template/1/2/3/4/5":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"passed": true,
				"message": "It's a hit!",`))
			case "/response/binary":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`"passed": true,
				"message": "It's a hit!",`))
			}
		case "POST":
			switch r.URL.Path {
			case "/form/string":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`4`))
			}
		}
	}))
}
