package https

import (
	"bytes"
	"fmt"
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

func FromFastHttpRequest(
	methodBytes []byte,
	url string,
	bodyBytes []byte,
	headerFunc func(setFunc func(k, v []byte)),
	requestURIBytes []byte,
	visitAllCookieFunc func(func(key, value []byte)),
) (*http.Request, error) {
	req, err := http.NewRequest(string(methodBytes), url, io.NopCloser(bytes.NewReader(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create http.Request: %w", err)
	}
	// Copy headers
	if headerFunc != nil {
		headerFunc(func(k, v []byte) { req.Header.Add(string(k), string(v)) })
	}
	// Set RequestURI (useful if mimicking server request)
	req.RequestURI = string(requestURIBytes)
	// Handle cookies explicitly
	if visitAllCookieFunc != nil {
		visitAllCookieFunc(func(k, v []byte) {
			req.AddCookie(&http.Cookie{Name: string(k), Value: string(v)})
		})
	}
	return req, nil
}
