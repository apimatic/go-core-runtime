package security

import (
	"testing"
)

// Assumptions (adjust/remove if your actual API differs):
//  - type DigestCodec interface { Encode([]byte) string; Decode(string) ([]byte, error) }
//  - func Hex() DigestCodec
//  - func Base64() DigestCodec
//  - func Base64URL() DigestCodec

func TestHexEncodeDecode(t *testing.T) {
	codec := NewHexDigestCodec()
	input := []byte{0x0A, 0x1B, 0xFF}
	encoded := codec.Encode(input)
	if encoded != "0a1bff" {
		t.Fatalf("expected 0a1bff, got %s", encoded)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if string(decoded) != string(input) {
		t.Fatalf("decoded mismatch: %v != %v", decoded, input)
	}
}

func TestHexEncodeEmpty(t *testing.T) {
	codec := NewHexDigestCodec()
	var input []byte
	encoded := codec.Encode(input)
	if encoded != "" {
		t.Fatalf("expected empty string, got %q", encoded)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected empty slice, got %v", decoded)
	}
}

func TestHexDecodeInvalidLength(t *testing.T) {
	codec := NewHexDigestCodec()
	if _, err := codec.Decode("abc"); err == nil {
		t.Fatalf("expected error for odd-length hex input")
	}
}

func TestHexDecodeInvalidCharacter(t *testing.T) {
	codec := NewHexDigestCodec()
	if _, err := codec.Decode("zzzzz"); err == nil {
		t.Fatalf("expected error for invalid hex characters")
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	codec := NewBase64DigestCodec()
	input := []byte{1, 2, 3, 4, 5}
	encoded := codec.Encode(input)
	if encoded != "AQIDBAU=" {
		t.Fatalf("expected AQIDBAU=, got %s", encoded)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if string(decoded) != string(input) {
		t.Fatalf("decoded mismatch: %v != %v", decoded, input)
	}
}

func TestBase64EncodeEmpty(t *testing.T) {
	codec := NewBase64DigestCodec()
	input := make([]byte, 0)
	encoded := codec.Encode(input)
	if encoded != "" {
		t.Fatalf("expected empty string, got %q", encoded)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected empty slice, got %v", decoded)
	}
}

func TestBase64DecodeInvalid(t *testing.T) {
	codec := NewBase64DigestCodec()
	if _, err := codec.Decode("!@#$"); err == nil {
		t.Fatalf("expected error for invalid base64 input")
	}
}

func TestBase64URLEncodeDecode(t *testing.T) {
	codec := NewBase64URLDigestCodec()
	input := []byte{10, 20, 30, 40, 50}
	encoded := codec.Encode(input)
	if encoded != "ChQeKDI" { // raw URL (no padding)
		t.Fatalf("expected ChQeKDI, got %s", encoded)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if string(decoded) != string(input) {
		t.Fatalf("decoded mismatch: %v != %v", decoded, input)
	}
}

func TestBase64URLEncodeEmpty(t *testing.T) {
	codec := NewBase64URLDigestCodec()
	var input []byte
	encoded := codec.Encode(input)
	if encoded != "" {
		t.Fatalf("expected empty string, got %q", encoded)
	}
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected empty slice, got %v", decoded)
	}
}

func TestBase64URLDecodeInvalid(t *testing.T) {
	codec := NewBase64URLDigestCodec()
	if _, err := codec.Decode("!@#$"); err == nil {
		t.Fatalf("expected error for invalid base64url input")
	}
}
