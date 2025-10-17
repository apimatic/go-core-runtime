package https

import (
	"bytes"
	"io"
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

func (ctx *HttpContext) GetResponseBody() ([]byte, error) {
	bodyBytes, err := io.ReadAll(ctx.Response.Body)
	ctx.Response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, err
}

func ReadRequestBody(req *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, err
}
