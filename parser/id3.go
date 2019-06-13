package parser

import (
	"errors"
	"fmt"
	"github.com/Monimena/audio-metadata"
	"github.com/mikkyang/id3-go"
	"github.com/mikkyang/id3-go/v1"
	"github.com/mikkyang/id3-go/v2"
	"io"
	"strconv"
	"strings"
)

var UnknownVersion = errors.New("error parsing metadata ID3: unknown version")

type ID3Parser struct {

}

func (p ID3Parser) Parse(file io.Reader) (*metadata.Info, error) {
	var info metadata.Info

	fSeeker, err := metadata.AsSeeker(file)

	if err != nil {
		return nil, err
	}

	if v2Tag := v2.ParseTag(fSeeker); v2Tag != nil {
		info =  mapMetadataID3(v2Tag)
	} else if v1Tag := v1.ParseTag(fSeeker); v1Tag != nil {
		info = mapMetadataID3(v1Tag)
	} else {
		return nil, UnknownVersion
	}

	return &info, nil
}

func mapMetadataID3(tagger id3.Tagger) metadata.Info {
	return metadata.Info{
		Title:   tagger.Title(),
		Artist:  tagger.Artist(),
		Album:   tagger.Album(),
		Year:    mapYear(tagger.Year()),
		Comment: mapComments(tagger.Comments()),
		Track:   0, // TODO: test this: The track number is stored in the last two bytes of the comment field. If the comment is 29 or 30 characters long, no track number can be stored.
		Genre:   tagger.Genre(),
		Other: map[string]string{
			"version": tagger.Version(),
			"size": fmt.Sprintf("%d", tagger.Size()),
			// TODO: add more fields? v2 frames
		},
	}
}

func mapYear(ystr string) int {
	year, err := strconv.Atoi(ystr)

	if err != nil {
		year = 0
	}

	return year
}

func mapComments(comments []string) string {
	return strings.Join(comments, " ")
}