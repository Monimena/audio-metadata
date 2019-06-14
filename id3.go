package metadata

import (
	"errors"
	"fmt"
	"github.com/mikkyang/id3-go"
	"github.com/mikkyang/id3-go/v1"
	"github.com/mikkyang/id3-go/v2"
	"io"
	"strconv"
	"strings"
)

var ErrUnknownVersion = errors.New("error parsing metadata ID3: unknown version")

type ID3Parser struct {

}

func (p ID3Parser) Parse(file io.Reader) (*Info, error) {
	var info Info

	fSeeker, err := asSeeker(file)

	if err != nil {
		return nil, err
	}

	if v2Tag := v2.ParseTag(fSeeker); v2Tag != nil {
		info =  fromID3Tagger(v2Tag)
	} else if v1Tag := v1.ParseTag(fSeeker); v1Tag != nil {
		info = fromID3Tagger(v1Tag)
	} else {
		return nil, ErrUnknownVersion
	}

	return &info, nil
}

func fromID3Tagger(tagger id3.Tagger) Info {
	return Info{
		Title:   tagger.Title(),
		Artist:  tagger.Artist(),
		Album:   tagger.Album(),
		Year:    mapYear(tagger.Year()),
		Comment: strings.Join(tagger.Comments(), " "), // TODO: check what these comments contain and why are they an array
		Track:   0, // TODO: test this: The track number is stored in the last two bytes of the comment field. If the comment is 29 or 30 characters long, no track number can be stored.
		Genre:   tagger.Genre(),
		Other: mapOther(tagger),
	}
}

func mapYear(ystr string) int {
	year, err := strconv.Atoi(ystr)

	if err != nil {
		year = 0
	}

	return year
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