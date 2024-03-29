package https

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// FileWrapper is a struct that represents a file along with its metadata such as the
// file content, file name, and file headers.
type FileWrapper struct {
	File        []byte
	FileName    string
	FileHeaders http.Header
}

// GetFile retrieves a file from the given fileUrl and returns it as a FileWrapper.
// It makes an HTTP GET request to the fileUrl to fetch the file's content and metadata.
func GetFile(fileUrl string) (FileWrapper, error) {
	url, err := url.Parse(fileUrl)
	if err != nil {
		return FileWrapper{}, internalError{Body: "Error parsing file", FileInfo: "fileWrapper.go/GetFile"}
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return FileWrapper{}, internalError{Body: "Error fetching file", FileInfo: "fileWrapper.go/GetFile"}
	}

	body, err := ReadBytes(resp.Body)

	file := FileWrapper{
		File:        body,
		FileName:    path.Base(url.Path),
		FileHeaders: resp.Header,
	}
	return file, err
}

// ReadBytes reads the data from the input io.Reader and returns it as a byte array.
// If there is an error while reading, it returns the error along with the byte array.
func ReadBytes(input io.Reader) ([]byte, error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		err = fmt.Errorf("error reading file: %v", err)
	}
	return bytes, err
}
