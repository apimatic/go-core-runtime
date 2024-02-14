package https

import (
	"net/http"
)

// HttpContext represents the HTTP request and response.
type HttpContext struct {
	Request  *http.Request
	Response *http.Response
}

func AddQuery(req *http.Request, key, value string) {
	queryVal := req.URL.Query()
	queryVal.Add(key, value)
	req.URL.RawQuery = encodeSpace(queryVal.Encode())
}