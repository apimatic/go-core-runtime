package https

import (
	"net/http"
)

// HttpContext represents the HTTP request and response.
type HttpContext struct {
	Request  *http.Request
	Response *http.Response
}
