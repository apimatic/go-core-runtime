package https

import (
	"testing"
)

func TestGetFile(t *testing.T) {
	file := GetFile("https://www.google.com/doodles/googles-new-logo")

	if file.FileName != "googles-new-logo" || len(file.File) <= 0 {
		t.Errorf("Expected Image File not recieved ")
	}
}

func TestGetFileErrorParsingUrl(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should panic because file url is empty.")
		}
	}()
	file := GetFile("")

	if file.FileName != "googles-new-logo" || len(file.File) <= 0 {
		t.Errorf("Expected Image File not recieved ")
	}
}

// func TestReadBytesInvalidResponse(t *testing.T) {
// 	// defer func() {
// 	// 	if r := recover(); r == nil {
// 	// 		t.Errorf("The code should panic because bytes to read are nil.")
// 	// 	}
// 	// }()
// 	input := io.NopCloser(strings.NewReader(""))
// 	//defer input.Close()
// 	bytes, err := ReadBytes(input)

// 	if len(bytes) <= 0 {
// 		t.Error(err)
// 	}
// }
