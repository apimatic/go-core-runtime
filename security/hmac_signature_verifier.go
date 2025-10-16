package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SignatureVerifier represents an interface for verifying the signature of an HTTP request.
type SignatureVerifier interface {
	// Verify checks the signature of the given HTTP request and returns
	// a VerificationResult indicating the outcome.
	Verify(req *http.Request) VerificationResult
}

// HmacSignatureVerifier verifies HMAC-based HTTP request signatures.
type HmacSignatureVerifier struct {
	// Name of the header carrying the provided signature (lookup is case-insensitive).
	signatureHeaderName string
	// HMAC algorithm used for signature generation (default: HmacSHA256).
	algorithm string
	// Initialized key spec; used to create a new Mac per verification call.
	secretKey []byte
	// Template containing "{digest}".
	signatureValueTemplate string
	// Resolves the bytes to sign from the request.
	requestBytesResolver func(http.Request) []byte
	// Codec used to decode (and possibly encode) digest text ↔ bytes (e.g., hex/base64).
	digestCodec DigestCodec
}

const signatureValuePlaceholder = "{digest}"

// NewHmacSignatureVerifier initializes a new HmacSignatureVerifier instance.
func NewHmacSignatureVerifier(
	secretKey string,
	signatureHeaderName string,
	digestCodec DigestCodec,
	requestBytesResolver func(http.Request) []byte,
	algorithm string,
	signatureValueTemplate string,
) (*HmacSignatureVerifier, error) {

	if strings.TrimSpace(secretKey) == "" {
		return nil, errors.New("secret key cannot be null or empty")
	}
	if strings.TrimSpace(signatureHeaderName) == "" {
		return nil, errors.New("signature header name cannot be null or empty")
	}
	if strings.TrimSpace(signatureValueTemplate) == "" {
		return nil, errors.New("signature value template cannot be null or empty")
	}
	if requestBytesResolver == nil {
		return nil, errors.New("requestBytesResolver cannot be null")
	}
	if digestCodec == nil {
		return nil, errors.New("digestCodec cannot be null")
	}
	if strings.TrimSpace(algorithm) == "" {
		return nil, errors.New("algorithm cannot be null or empty")
	}

	return &HmacSignatureVerifier{
		signatureHeaderName:    signatureHeaderName,
		algorithm:              algorithm,
		secretKey:              []byte(secretKey),
		signatureValueTemplate: signatureValueTemplate,
		requestBytesResolver:   requestBytesResolver,
		digestCodec:            digestCodec,
	}, nil
}

// Verify verifies the HMAC signature for the given request.
func (v *HmacSignatureVerifier) Verify(request http.Request) VerificationResult {
	headerValue := []string{""}
	for k, val := range request.Header {
		if strings.EqualFold(k, v.signatureHeaderName) {
			headerValue = val
			break
		}
	}

	if len(headerValue) <= 1 {
		return NewFailure(fmt.Sprintf("Signature header '%s' is missing.", v.signatureHeaderName))
	}

	provided := v.extractSignature(headerValue[0])
	if len(provided) == 0 {
		return NewFailure(fmt.Sprintf("Malformed signature header '%s'.", v.signatureHeaderName))
	}

	message := v.requestBytesResolver(request)
	computed, err := v.computeHMAC(message)
	if err != nil {
		return NewFailure(fmt.Sprintf("Error computing HMAC: %v", err))
	}

	if subtle.ConstantTimeCompare(provided, computed) == 1 {
		return NewSuccess()
	}
	return NewFailure("Signature verification failed.")
}

// computeHMAC computes the HMAC for the given message using the configured algorithm.
func (v *HmacSignatureVerifier) computeHMAC(message []byte) ([]byte, error) {
	switch strings.ToUpper(v.algorithm) {
	case "HMACSHA256", "HMAC-SHA256", "SHA256":
		mac := hmac.New(sha256.New, v.secretKey)
		mac.Write(message)
		return mac.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", v.algorithm)
	}
}

// extractSignature extracts and decodes the digest from the header value.
func (v *HmacSignatureVerifier) extractSignature(headerValue string) []byte {
	index := strings.Index(v.signatureValueTemplate, signatureValuePlaceholder)
	if index < 0 {
		return []byte{}
	}

	prefix := v.signatureValueTemplate[:index]
	suffix := v.signatureValueTemplate[index+len(signatureValuePlaceholder):]

	prefixAt := indexOfIgnoreCase(headerValue, prefix, 0)
	if prefixAt < 0 {
		return []byte{}
	}

	digestStart := prefixAt + len(prefix)
	digestEnd := len(headerValue)

	if suffix != "" {
		suffixAt := indexOfIgnoreCase(headerValue, suffix, digestStart)
		if suffixAt < 0 {
			return []byte{}
		}
		digestEnd = suffixAt
	}

	if digestEnd < digestStart {
		return []byte{}
	}

	digest := strings.TrimSpace(headerValue[digestStart:digestEnd])

	// Strip optional quotes
	if len(digest) >= 2 && digest[0] == '"' && digest[len(digest)-1] == '"' {
		digest = digest[1 : len(digest)-1]
	}

	defer func() {
		if r := recover(); r != nil {
			_ = r // ignore panic from invalid hex/base64
		}
	}()

	decoded, _ := v.digestCodec.Decode(digest)
	return decoded
}

// indexOfIgnoreCase finds the first occurrence of needle in haystack (case-insensitive).
func indexOfIgnoreCase(haystack, needle string, fromIndex int) int {
	if needle == "" {
		return fromIndex
	}
	haystackLower := strings.ToLower(haystack)
	needleLower := strings.ToLower(needle)
	idx := strings.Index(haystackLower[fromIndex:], needleLower)
	if idx == -1 {
		return -1
	}
	return fromIndex + idx
}
