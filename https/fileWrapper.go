package https

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type FileWrapper struct {
	File        []byte
	FileName    string
	FileHeaders http.Header
}

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

func ReadBytes(input io.Reader) ([]byte, error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		err = fmt.Errorf("Error reading file: %v", err)
	}
	return bytes, err
}
