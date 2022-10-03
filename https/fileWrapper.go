package https

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

type FileWrapper struct {
	File        []byte
	FileName    string
	FileHeaders http.Header
}

func GetFile(fileUrl string) FileWrapper {
	url, _ := url.Parse(fileUrl)
	resp, err := http.Get(url.String())
	if err != nil {
		log.Panic(err)
	}

	body, err := ReadBytes(resp.Body)

	file := FileWrapper{
		File:        body,
		FileName:    path.Base(url.Path),
		FileHeaders: resp.Header,
	}
	return file
}

func ReadBytes(input io.Reader) ([]byte, error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		log.Panic(err)
	}
	return bytes, err
}
