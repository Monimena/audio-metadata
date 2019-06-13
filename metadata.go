package metadata

import (
	"bytes"
	"errors"
	"github.com/Monimena/audio-metadata/parser"
	"io"
	"io/ioutil"
	"net/http"
)

var UnknownContentType = errors.New("error parsing metadata: unknown content type")

var id3Parser = parser.ID3Parser{}
var vorbisParser = parser.VorbisParser{}

type Parser interface {
	Parse(file io.Reader) (*Info, error)
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
	"audio/mpeg": &id3Parser,
	"audio/MPA": &id3Parser,
	"audio/mpa-robust": &id3Parser,
	"audio/vnd.wave": &id3Parser,
	"audio/wav": &id3Parser,
	"audio/wave": &id3Parser,
	"audio/x-wav": &id3Parser,
	"audio/x-aiff": &id3Parser,
	"audio/aiff": &id3Parser,

	"audio/ogg": &vorbisParser,
	"audio/opus": &vorbisParser,
	"audio/flac": &vorbisParser,
	"audio/vorbis": &vorbisParser,
}

func Parse(file io.Reader) (*Info, error) {
	fSeeker, err := AsSeeker(file)

	if err != nil {
		return nil, err
	}

	b := make([]byte, 512)
	_, err = fSeeker.Read(b)

	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(b)

	if p, found := parserMap[contentType]; found {
		return p.Parse(fSeeker)
	}

	return nil, UnknownContentType
}

func AsSeeker(r io.Reader) (io.ReadSeeker, error) {
	if rs, ok := r.(io.ReadSeeker); ok {
		_, err := rs.Seek(0, io.SeekStart) // reset

		if err != nil {
			return nil, err
		}

		return rs, nil                  // r is already a readSeeker under the hood, return it
	}

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}