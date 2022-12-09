package testHelper

import (
	"net/http"
	"testing"
)

func TestCheckResponseHeaders(t *testing.T) {
	result := http.Header{
		"Content-Type":    {"application/responseType"},
		"Accept":          {"application/noTerm"},
		"Accept-Encoding": {"UTF-8"},
	}
	expected := []TestHeader{
		NewTestHeader(true, "Content-Type", "application/responseType"),
		NewTestHeader(true, "Accept", "application/noTerm"),
		NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	CheckResponseHeaders(t, result, expected, false)
}

func TestCheckResponseHeadersError(t *testing.T) {
	result := http.Header{
		"Host":            {"www.host.com"},
		"Content-Type":    {""},
		"Authorization":   {"Bearer Token"},
		"Accept":          {"application/noTerm"},
		"Accept-Encoding": {"UTF-8"},
	}
	expected := []TestHeader{
		NewTestHeader(true, "Content-Type", "application/responseType"),
		NewTestHeader(true, "Accept", "application/noTerm"),
		NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	CheckResponseHeaders(&testing.T{}, result, expected, false)
}

func TestCheckResponseHeadersAllowExtras(t *testing.T) {
	result := http.Header{
		"Host":            {"www.host.com"},
		"Content-Type":    {"application/responseType"},
		"Accept":          {"application/noTerm"},
		"Accept-Encoding": {"UTF-8"},
		"Authorization":   {"Bearer Token"},
	}
	expected := []TestHeader{
		NewTestHeader(true, "Content-Type", "application/responseType"),
		NewTestHeader(true, "Accept", "application/noTerm"),
		NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	CheckResponseHeaders(t, result, expected, true)
}

func TestCheckResponseHeadersAllowExtrasError(t *testing.T) {
	result := http.Header{
		"Host":          {"www.host.com"},
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer Token"},
	}
	expected := []TestHeader{
		NewTestHeader(true, "Content-Type", "application/responseType"),
		NewTestHeader(true, "Accept", "application/noTerm"),
		NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	CheckResponseHeaders(&testing.T{}, result, expected, true)
}
