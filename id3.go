package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mikkyang/id3-go"
	"github.com/mikkyang/id3-go/v1"
	"github.com/mikkyang/id3-go/v2"
	"io"
	"log"
	"strconv"
	"strings"
)

var ErrUnknownVersion = errors.New("error parsing metadata ID3: unknown version\n")

func ParseID3(file io.Reader) (*Info, error) {
	var info Info

	fSeeker, err := asSeeker(file)

	if err != nil {
		return nil, err
	}

	if v2Tag := v2.ParseTag(fSeeker); v2Tag != nil {
		info = fromID3Tagger(v2Tag, 2)
	} else if v1Tag := v1.ParseTag(fSeeker); v1Tag != nil {
		info = fromID3Tagger(v1Tag, 1)
	} else {
		return nil, ErrUnknownVersion
	}

	return &info, nil
}

func fromID3Tagger(tagger id3.Tagger, version int) Info {
	return Info{
		Title:   trimNullChar(tagger.Title()),
		Artist:  trimNullChar(tagger.Artist()),
		Album:   trimNullChar(tagger.Album()),
		Year:    mapYear(trimNullChar(tagger.Year())),
		Comment: trimNullChar(strings.Join(tagger.Comments(), "\n")),
		Track:   mapTrack(tagger, version),
		Genre:   trimNullChar(tagger.Genre()),
		Other:   mapOther(tagger),
	}
}

func trimNullChar(s string) string {
	return strings.TrimSuffix(s, "\000")
}

func mapYear(ystr string) int {
	year, err := strconv.Atoi(ystr)

	if err != nil {
		year = 0
	}

	return year
}

func mapTrack(tagger id3.Tagger, version int) int {
	var t int
	var err error

	if version == 2 {
		// ID3v2 has the tracknumber as a frames field
		log.Println("ID3v2")

		t, err = strconv.Atoi(strings.Split(tagger.Frame("TRCK").String(), "/")[0])
	} else if version == 1 {
		// ID3v1 has the tracknumber in the last byte of the comments
		log.Println("ID3v1")
		comments := trimNullChar(tagger.Comments()[0])

		l := len(comments)

		log.Printf("comment length: %d\n", l)

		if bytes.Compare([]byte(comments[l-1:l-1]), make([]byte, 0)) == 0 {
			t, err = strconv.Atoi(comments[l:l])
		}
	} else {
		// unknown version
		return 0
	}

	if err != nil {
		return 0
	}

	return t
}

func mapOther(tagger id3.Tagger) map[string]string {
	// TODO: add more fields? v2 frames
	var other = map[string]string{}

	if len(tagger.Version()) > 0 {
		other["version"] = tagger.Version()
	}

	if tagger.Size() > 0 {
		other["size"] = fmt.Sprintf("%d", tagger.Size())
	}

	return other
}
