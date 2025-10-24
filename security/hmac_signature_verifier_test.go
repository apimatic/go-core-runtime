package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/apimatic/go-core-runtime/https"
)

const (
	SECRET_KEY               = "testSecret"
	SIGNATURE_HEADER         = "X-Signature"
	BODY                     = "payload"
	DEFAULT_ALGORITHM        = "HmacSHA256"
	SIGNATURE_VALUE_TEMPLATE = "{digest}"
)

var DEFAULT_DIGEST_CODEC = NewHexDigestCodec()

// Helper to compute expected signature (HmacSHA256 only for these tests).
func computeHmacHex(key, data string, algorithm func() hash.Hash, codec DigestCodec) string {
	m := hmac.New(algorithm, []byte(key))
	m.Write([]byte(data))
	return codec.Encode(m.Sum(nil))
}

// --- Constructor validation tests ---
func TestHmacVerifierConstructorValidation(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		header  string
		codec   DigestCodec
		res     func(*http.Request) ([]byte, error)
		algo    string
		tmpl    string
		wantErr bool
	}{
		{"NilSecret", "", SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, true},
		{"BlankSecret", "   ", SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, true},
		{"NilHeader", "k", "", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, true},
		{"BlankHeader", "k", "   ", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, true},
		{"NilCodec", "k", "X", nil, https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, true},
		{"Nilhttps.ReadRequestBody", "k", "X", DEFAULT_DIGEST_CODEC, nil, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, true},
		{"NilAlgo", "k", "X", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, "", SIGNATURE_VALUE_TEMPLATE, true},
		{"BlankAlgo", "k", "X", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, "   ", SIGNATURE_VALUE_TEMPLATE, true},
		{"NilTemplate", "k", "X", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, "", true},
		{"BlankTemplate", "k", "X", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, "   ", true},
		{"Valid", "k", "X", DEFAULT_DIGEST_CODEC, https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewHmacSignatureVerifier(tc.secret, tc.header, tc.codec, tc.res, tc.algo, tc.tmpl)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func newRequest(method string, headers map[string][]string, body string) *http.Request {
	return &http.Request{
		Method:        method,
		Header:        headers,
		Body:          io.NopCloser(bytes.NewBuffer([]byte(body))),
		ContentLength: int64(len(body)),
	}
}

// --- Behaviour tests ---
func TestHmacVerifySuccess(t *testing.T) {
	signature := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {signature}}, BODY)
	verifier, _ := NewHmacSignatureVerifier(
		SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC, https.ReadRequestBody,
		DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if !result.Success() {
		t.Fatalf("expected success, got errors: %v", result.Errors())
	}
}

func TestHmacVerifySuccessSHA512(t *testing.T) {
	signature := computeHmacHex(SECRET_KEY, BODY, sha512.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {signature}}, BODY)
	verifier, _ := NewHmacSignatureVerifier(
		SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC, https.ReadRequestBody,
		"HmacSHA512", SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if !result.Success() {
		t.Fatalf("expected success, got errors: %v", result.Errors())
	}
}

func TestHmacVerifySuccessSHA520(t *testing.T) {
	signature := computeHmacHex(SECRET_KEY, BODY, sha512.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {signature}}, BODY)
	verifier, _ := NewHmacSignatureVerifier(
		SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC, https.ReadRequestBody,
		"HmacSHA520", SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if result.Success() {
		t.Fatalf("expected failure")
	}
}

func TestHmacVerifyMissingHeader(t *testing.T) {
	req := newRequest("POST", map[string][]string{}, BODY)
	verifier, _ := NewHmacSignatureVerifier(
		SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if result.Success() {
		t.Fatalf("expected failure")
	}
	if len(result.Errors()) == 0 || !strings.Contains(result.Errors()[0], "missing") {
		t.Fatalf("unexpected error: %v", result.Errors())
	}
}

func TestHmacVerifyMalformedSignature(t *testing.T) {
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {"not-a-hex-signature"}}, BODY)
	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if result.Success() {
		t.Fatalf("expected failure")
	}
	if len(result.Errors()) == 0 || !strings.Contains(strings.ToLower(result.Errors()[0]), "malformed") {
		t.Fatalf("unexpected error: %v", result.Errors())
	}
}

func TestHmacVerifyWrongSignature(t *testing.T) {
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {"deadbeef"}}, BODY)
	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if result.Success() {
		t.Fatalf("expected failure")
	}
}

func TestHmacVerifyResolverException(t *testing.T) {
	sig := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {sig}}, BODY)
	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		func(r *http.Request) ([]byte, error) { return nil, errors.New("resolver error") },
		DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)

	result := verifier.Verify(req)
	if result.Success() {
		t.Fatalf("expected failure")
	}
	if len(result.Errors()) == 0 || !strings.Contains(result.Errors()[0], "resolver error") {
		t.Fatalf("unexpected errors: %v", result.Errors())
	}
}

// --- Template variation tests ---

func TestHmacVerifyTemplatePrefixSHA512(t *testing.T) {
	digest := computeHmacHex(SECRET_KEY, BODY, sha512.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {"sha256=" + digest}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		"HmacSHA512", "sha256={digest}")

	if sv := verifier.Verify(req); !sv.Success() {
		t.Fatalf("expected success")
	}
}

func TestHmacVerifyTemplatePrefix(t *testing.T) {
	digest := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {"sha256=" + digest}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, "sha256={digest}")

	if sv := verifier.Verify(req); !sv.Success() {
		t.Fatalf("expected success")
	}
}

func TestHmacVerifyTemplatePrefixWrong(t *testing.T) {
	digest := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {digest}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, "sha256={digest}")

	sv := verifier.Verify(req)
	if sv.Success() {
		t.Fatalf("expected failure")
	}

	if len(sv.Errors()) == 0 || sv.Errors()[0] != "Malformed signature header '"+SIGNATURE_HEADER+"'." {
		t.Fatalf("unexpected errors: %v", sv.Errors())
	}
}

func TestHmacVerifyTemplateSuffix(t *testing.T) {
	digest := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {digest + "complex"}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, "{digest}complex")

	if sv := verifier.Verify(req); !sv.Success() {
		t.Fatalf("expected success")
	}
}

func TestHmacVerifyTemplateSuffixWrong(t *testing.T) {
	digest := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {digest}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, "{digest}complex")

	sv := verifier.Verify(req)
	if sv.Success() {
		t.Fatalf("expected failure")
	}

	if len(sv.Errors()) == 0 || sv.Errors()[0] != "Malformed signature header '"+SIGNATURE_HEADER+"'." {
		t.Fatalf("unexpected errors: %v", sv.Errors())
	}
}

func TestHmacVerifyTemplateMismatchPlaceholder(t *testing.T) {
	d := computeHmacHex(SECRET_KEY, BODY, sha256.New, DEFAULT_DIGEST_CODEC)
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {d}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, "{diet}") // wrong placeholder

	if sv := verifier.Verify(req); sv.Success() {
		t.Fatalf("expected failure due to malformed template")
	}
}

func TestHmacVerifyTemplateQuotedWrong(t *testing.T) {
	req := newRequest("POST", map[string][]string{SIGNATURE_HEADER: {"v0=\"a\"complex"}}, BODY)

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, DEFAULT_DIGEST_CODEC,
		https.ReadRequestBody,
		DEFAULT_ALGORITHM, "v0={digest}complex")

	if sv := verifier.Verify(req); sv.Success() {
		t.Fatalf("expected failure")
	}
}

// --- Complex resolver tests (base64 codec variation) ---

func TestHmacVerifyComplexResolverBody(t *testing.T) {
	codec := NewBase64DigestCodec()
	body := `{"id":123,"type":"payment","amount":100.5}`
	// Compose data string replicating Java logic (extracting "payment" as /type)
	data := "session=abc123:2025-09-17T12:34:56Z:POST:payment" +
		":{\"id\":123,\"type\":\"payment\",\"amount\":100.5}"

	signature := computeHmacHex(SECRET_KEY, data, sha256.New, codec)
	headers := map[string][]string{
		SIGNATURE_HEADER: {signature},
		"Cookie":         {"session=abc123"},
		"X-Timestamp":    {"2025-09-17T12:34:56Z"},
	}
	req := newRequest("POST", headers, body)

	resolver := func(r *http.Request) ([]byte, error) {
		cookie := r.Header.Get("Cookie")
		timestamp := r.Header.Get("X-Timestamp")
		// naive extraction of "type" value from JSON (test-only)
		var extracted string
		reqBody, _ := https.ReadRequestBody(r)
		if idx := strings.Index(string(reqBody), `"type":"`); idx >= 0 {
			sub := reqBody[idx+len(`"type":"`):]
			if j := bytes.IndexByte(sub, '"'); j >= 0 {
				extracted = string(sub[:j])
			}
		}
		payload := strings.Join([]string{
			cookie, timestamp, r.Method, extracted, string(reqBody),
		}, ":")
		return []byte(payload), nil
	}

	verifier, _ := NewHmacSignatureVerifier(SECRET_KEY, SIGNATURE_HEADER, codec, resolver, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)
	if sv := verifier.Verify(req); !sv.Success() {
		t.Fatalf("expected success")
	}
}

func TestHmacVerifyComplexResolverHeaders(t *testing.T) {
	const secret = SECRET_KEY
	header := SIGNATURE_HEADER
	codec := NewBase64DigestCodec()
	body := `{"id":123,"type":"payment","amount":100.5}`

	data := "session=abc123:2025-09-17T12:34:56Z:POST:x-signature-header-value" +
		":{\"id\":123,\"type\":\"payment\",\"amount\":100.5}"

	signature := computeHmacHex(secret, data, sha256.New, codec)
	headers := map[string][]string{
		header:           {signature},
		"Cookie":         {"session=abc123"},
		"X-Timestamp":    {"2025-09-17T12:34:56Z"},
		SIGNATURE_HEADER: {"x-signature-header-value"},
	}
	req := newRequest("POST", headers, body)

	resolver := func(r *http.Request) ([]byte, error) {
		reqBody, _ := https.ReadRequestBody(r)
		payload := strings.Join([]string{
			r.Header.Get("Cookie"),
			r.Header.Get("X-Timestamp"),
			r.Method,
			string(reqBody),
		}, ":")
		return []byte(payload), nil
	}

	verifier, _ := NewHmacSignatureVerifier(secret, header, codec, resolver, DEFAULT_ALGORITHM, SIGNATURE_VALUE_TEMPLATE)
	if sv := verifier.Verify(req); sv.Success() {
		t.Fatalf("expected failure due to wrong signature (header value not included)")
	}
}
