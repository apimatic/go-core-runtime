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

func TestGetFileErrorParsingUrl(t *testing.T) {
	_, err := GetFile("")
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

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestReadBytesInvalidResponse(t *testing.T) {
	bytes, err := ReadBytes(errReader(0))

	if len(bytes) != 0 {
		t.Error(err)
	}
}
