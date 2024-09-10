package https

import (
	"errors"
	"testing"
)

func TestGetFile(t *testing.T) {
	file, err := GetFile("https://www.google.com/doodles/googles-new-logo")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}

	if file.FileName != "googles-new-logo" || len(file.File) <= 0 {
		t.Errorf("Expected Image File not recieved ")
	}
}

func TestGetFileWithContentType(t *testing.T) {
	file, err := GetFileWithContentType("https://www.google.com/doodles/googles-new-logo", "image/png")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}

	if file.FileName != "googles-new-logo" ||
		file.FileHeaders.Get(CONTENT_TYPE_HEADER) != "image/png" ||
		len(file.File) <= 0 {
		t.Errorf("Expected Image File not recieved ")
	}
}

func TestGetFileFromLocalPath(t *testing.T) {
	file, err := GetFile("../internal/binary.png")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}

	if file.FileName != "binary.png" ||
		file.FileHeaders.Get(CONTENT_TYPE_HEADER) != OCTET_STREAM_CONTENT_TYPE ||
		len(file.File) <= 0 {
		t.Errorf("Expected Image File not recieved ")
	}
}

func TestGetFileWithContentTypeFromLocalPath(t *testing.T) {
	file, err := GetFileWithContentType("../internal/binary.png", "image/png")
	if err != nil {
		t.Errorf("GetFile failed: %v", err)
	}

	if file.FileName != "binary.png" ||
		file.FileHeaders.Get(CONTENT_TYPE_HEADER) != "image/png" ||
		len(file.File) <= 0 {
		t.Errorf("Expected Image File not recieved ")
	}
}

func TestGetFileErrorParsingUrl(t *testing.T) {
	_, err := GetFile("")
	if err == nil {
		t.Errorf("GetFile failed: %v", err)
	}
}

func TestGetFileWithContentTypeErrorParsingUrl(t *testing.T) {
	_, err := GetFileWithContentType("", "image/png")
	if err == nil {
		t.Errorf("GetFile failed: %v", err)
	}
}

func TestGetFileErrorParsingUrlWithSpecialChar(t *testing.T) {
	_, err := GetFile("hhhh%#")
	if err == nil {
		t.Errorf("GetFile failed: %v", err)
	}
}

func TestGetFileWithContentTypeErrorParsingUrlWithSpecialChar(t *testing.T) {
	_, err := GetFileWithContentType("hhhh%#", "image/png")
	if err == nil {
		t.Errorf("GetFile failed: %v", err)
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestReadBytesInvalidResponse(t *testing.T) {
	bytes, err := readBytes(errReader(0))

	if len(bytes) != 0 {
		t.Error(err)
	}
}
