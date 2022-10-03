package https

import (
	"net/http"
)

type HttpContext struct {
	Request  *http.Request
	Response *http.Response
}
