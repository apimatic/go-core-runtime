package https

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// FileWrapper is a struct that represents a file along with its metadata such as the
// file content, file name, and file headers.
type FileWrapper struct {
	File        []byte
	FileName    string
	FileHeaders http.Header
}

func (f FileWrapper) String() string {
	return fmt.Sprintf("FileWrapper[FileName=%v]", f.FileName)
}

// isURL checks if the given parsedPath is a URL
func isURL(parsedPath *url.URL) bool {
	return parsedPath.Scheme == "http" || parsedPath.Scheme == "https"
}

// GetFile retrieves a file from the given filePath and returns it as a FileWrapper.
// It makes an HTTP GET request to the filePath to fetch the file's content and metadata.
// OR It uses os.ReadFile to read the file's content and metadata.
func GetFile(filePath string) (FileWrapper, error) {
	_fileWrapper := FileWrapper{FileHeaders: http.Header{}}
	parsedPath, err := url.Parse(filePath)
	if err != nil {
		return _fileWrapper, internalError{Body: "Error parsing file", FileInfo: "fileWrapper.go/GetFile"}
	}

	if isURL(parsedPath) {
		resp, err := http.Get(parsedPath.String())
		if err != nil {
			return _fileWrapper, internalError{Body: "Error fetching file", FileInfo: "fileWrapper.go/GetFile"}
		}
		_fileWrapper.FileName = path.Base(parsedPath.Path)
		_fileWrapper.FileHeaders = resp.Header
		_fileWrapper.File, err = readBytes(resp.Body)
		return _fileWrapper, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return _fileWrapper, err
	}
	_fileWrapper.File, err = os.ReadFile(filePath)
	if err != nil {
		return _fileWrapper, internalError{Body: "Error reading file", FileInfo: "fileWrapper.go/GetFile"}
	}
	_fileWrapper.FileName = filepath.Base(filePath)
	_fileWrapper.FileHeaders.Set(CONTENT_TYPE_HEADER, OCTET_STREAM_CONTENT_TYPE)
	return _fileWrapper, err
}

// GetFileWithContentType retrieves a file from the given filePath using GetFile and returns it as a FileWrapper.
// It also sets the provided "content-type" in the file headers.
func GetFileWithContentType(filePath string, contentType string) (FileWrapper, error) {
	_fileWrapper, err := GetFile(filePath)
	_fileWrapper.FileHeaders.Set(CONTENT_TYPE_HEADER, contentType)
	return _fileWrapper, err
}

// readBytes reads the data from the input io.Reader and returns it as a byte array.
// If there is an error while reading, it returns the error along with the byte array.
func readBytes(input io.Reader) ([]byte, error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		err = fmt.Errorf("error reading file: %v", err)
	}
	return bytes, err
}
