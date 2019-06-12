package metadata

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Parser interface {
	Parse(file io.Reader) (Info, error)
}

type Info struct {
	Title   string
	Artist  string
	Album   string
	Year    int
	Comment string
	Track   int
	Genre   string
	Other   map[string]string
}

var parserMap = map[string]Parser{

}

func Parse(file io.Reader) (Info, error) {
	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	contentType := http.DetectContentType(fileBytes)

	return parserMap[contentType].Parse(bytes.NewReader(fileBytes))
}