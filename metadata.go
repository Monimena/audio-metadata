package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var ErrUnknownContentType = errors.New("error parsing metadata: unknown content type\n")

type Parser interface {
	Parse(file io.Reader) (*Info, error)
}

type ParserFunc func(file io.Reader) (*Info, error)

func (pf ParserFunc) Parse(file io.Reader) (*Info, error) {
	return pf(file)
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

type MIMEParser struct {
	parserMap map[string]Parser
}

type FallbackParser struct {
	parsers []Parser
} // TODO: add parse method that tries to parse with every parser else returns error
// TODO: add append method

var defaultMIMEParser = MIMEParser{
	parserMap: map[string]Parser{
		"audio/mpeg":       ParserFunc(ParseID3),
		"audio/MPA":        ParserFunc(ParseID3),
		"audio/mpa-robust": ParserFunc(ParseID3),
		"audio/vnd.wave":   ParserFunc(ParseID3),
		"audio/wav":        ParserFunc(ParseID3),
		"audio/wave":       ParserFunc(ParseID3),
		"audio/x-wav":      ParserFunc(ParseID3),
		"audio/x-aiff":     ParserFunc(ParseID3),
		"audio/aiff":       ParserFunc(ParseID3),

		"audio/ogg":    ParserFunc(ParseVorbis),
		"audio/opus":   ParserFunc(ParseVorbis),
		"audio/flac":   ParserFunc(ParseVorbis),
		"audio/vorbis": ParserFunc(ParseVorbis),
	},
}

var fallbackParser = FallbackParser{
	parsers: []Parser{&defaultMIMEParser, ParserFunc(ParseID3), ParserFunc(ParseVorbis)},
}

func (mp *MIMEParser) Parse(file io.Reader) (*Info, error) {
	fSeeker, err := asSeeker(file)

	if err != nil {
		return nil, err
	}

	b := make([]byte, 512)
	_, err = fSeeker.Read(b)

	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(b)

	fmt.Printf("content type detected: %s\n", contentType)

	if p, found := mp.parserMap[contentType]; found {
		return p.Parse(fSeeker)
	}

	return nil, ErrUnknownContentType
}

func (mp *MIMEParser) Append(s string, p Parser) {
	if mp.parserMap == nil {
		mp.parserMap = map[string]Parser{}
	}

	mp.parserMap[s] = p
}

func Parse(file io.Reader) (*Info, error) {
	return defaultMIMEParser.Parse(file)
}

func asSeeker(r io.Reader) (io.ReadSeeker, error) {
	if rs, ok := r.(io.ReadSeeker); ok {
		_, err := rs.Seek(0, io.SeekStart) // reset

		if err != nil {
			return nil, err
		}

		return rs, nil // r is already a readSeeker under the hood, return it
	}

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
