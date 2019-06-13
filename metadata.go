package metadata

import (
	"bytes"
	"errors"
	"github.com/Monimena/audio-metadata/parser"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const unknownContentType = "error parsing metadata: unknown content type"

type Parser interface {
	Parse(file io.ReadSeeker) (Info, error)
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
	"audio/mpeg": parser.ID3Parser{},
}

func Parse(file io.Reader) (Info, error) {
	var info Info
	var err error


	fSeeker := asSeeker(file)

	b := make([]byte, 512)
	_, err = fSeeker.Read(b)
	fSeeker.Seek(0, io.SeekStart) // reset

	if err != nil {
		return info, err
	}

	contentType := http.DetectContentType(b)

	if parserMap[contentType] == nil {
		return info, errors.New(unknownContentType)
	}

	return parserMap[contentType].Parse(fSeeker)
}

func asSeeker(r io.Reader) io.ReadSeeker {
	if rs, ok := r.(io.ReadSeeker); ok {
		rs.Seek(0, io.SeekStart) // reset
		return rs // r is already a readSeeker under the hood, return it
	}

	b, err := ioutil.ReadAll(r)

	if err != nil {
		log.Fatal(err)
	}

	return bytes.NewReader(b)
}