package testHelper_test

import (
	"github.com/apimatic/go-core-runtime/testHelper"
	"net/http"
	"testing"
)

func TestCheckResponseHeaders(t *testing.T) {
	result := http.Header{
		"Content-Type":    {"application/responseType"},
		"Accept":          {"application/noTerm"},
		"Accept-Encoding": {"UTF-8"},
	}
	expected := []testHelper.TestHeader{
		testHelper.NewTestHeader(true, "Content-Type", "application/responseType"),
		testHelper.NewTestHeader(true, "Accept", "application/noTerm"),
		testHelper.NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	testHelper.CheckResponseHeaders(t, result, expected, false)
}

func TestCheckResponseHeadersError(t *testing.T) {
	result := http.Header{
		"Host":            {"www.host.com"},
		"Content-Type":    {""},
		"Authorization":   {"Bearer Token"},
		"Accept":          {"application/noTerm"},
		"Accept-Encoding": {"UTF-8"},
	}
	expected := []testHelper.TestHeader{
		testHelper.NewTestHeader(true, "Content-Type", "application/responseType"),
		testHelper.NewTestHeader(true, "Accept", "application/noTerm"),
		testHelper.NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	testHelper.CheckResponseHeaders(&testing.T{}, result, expected, false)
}

func TestCheckResponseHeadersAllowExtras(t *testing.T) {
	result := http.Header{
		"Host":            {"www.host.com"},
		"Content-Type":    {"application/responseType"},
		"Accept":          {"application/noTerm"},
		"Accept-Encoding": {"UTF-8"},
		"Authorization":   {"Bearer Token"},
	}
	expected := []testHelper.TestHeader{
		testHelper.NewTestHeader(true, "Content-Type", "application/responseType"),
		testHelper.NewTestHeader(true, "Accept", "application/noTerm"),
		testHelper.NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	testHelper.CheckResponseHeaders(t, result, expected, true)
}

func TestCheckResponseHeadersAllowExtrasError(t *testing.T) {
	result := http.Header{
		"Host":          {"www.host.com"},
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer Token"},
	}
	expected := []testHelper.TestHeader{
		testHelper.NewTestHeader(true, "Content-Type", "application/responseType"),
		testHelper.NewTestHeader(true, "Accept", "application/noTerm"),
		testHelper.NewTestHeader(true, "Accept-Encoding", "UTF-8"),
	}
	testHelper.CheckResponseHeaders(&testing.T{}, result, expected, true)
}
