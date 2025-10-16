package security

import (
	"encoding/base64"
	"encoding/hex"
)

// DigestCodec is an interface for encoding and decoding digest values.
type DigestCodec interface {
	// Encode encodes a byte slice digest into a string representation.
	Encode(bytes []byte) string

	// Decode decodes a string representation back into a byte slice.
	Decode(encoded string) ([]byte, error)
}

// NewHexDigestCodec returns a DigestCodec for Hex encoding/decoding.
func NewHexDigestCodec() DigestCodec {
	return hexDigestCodec{}
}

// NewBase64DigestCodec returns a DigestCodec for Base64 encoding/decoding.
func NewBase64DigestCodec() DigestCodec {
	return base64DigestCodec{}
}

// NewBase64URLDigestCodec returns a DigestCodec for Base64 URL-safe encoding/decoding.
func NewBase64URLDigestCodec() DigestCodec {
	return base64URLDigestCodec{}
}

// hexDigestCodec implements DigestCodec for hexadecimal encoding.
type hexDigestCodec struct{}

func (hexDigestCodec) Encode(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func (hexDigestCodec) Decode(encoded string) ([]byte, error) {
	return hex.DecodeString(encoded)
}

// base64DigestCodec implements DigestCodec for Base64 encoding.
type base64DigestCodec struct{}

func (base64DigestCodec) Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func (base64DigestCodec) Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// base64URLDigestCodec implements DigestCodec for Base64 URL-safe encoding.
type base64URLDigestCodec struct{}

func (base64URLDigestCodec) Encode(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func (base64URLDigestCodec) Decode(encoded string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(encoded)
}
