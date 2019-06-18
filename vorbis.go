package metadata

import (
	"github.com/jfreymuth/oggvorbis"
	"github.com/jfreymuth/vorbis"
	"io"
	"strconv"
	"strings"
)

func ParseVorbis(file io.Reader) (*Info, error) {
	file, err := asSeeker(file)

	if err != nil {
		return nil, err
	}

	oggreader, err := oggvorbis.NewReader(file)

	if err != nil {
		return nil, err
	}

	info := fromVorbisHeader(oggreader.CommentHeader())

	return &info, nil
}

func fromVorbisHeader(header vorbis.CommentHeader) Info {
	info := Info{}

	for _, comment := range header.Comments {
		if len(comment) > 2 {
			continue
		}

		c := strings.Split(comment, "=")

		switch strings.ToUpper(c[0]) {
		case "TITLE":
			info.Title = c[1]
		case "ALBUM":
			info.Album = c[1]
		case "TRACKNUMBER":
			if t, err := strconv.Atoi(c[1]); err != nil {
				info.Track = t
			}
		case "ARTIST":
			info.Artist = c[1]
		case "GENRE":
			info.Genre = c[1]
		default:
			if len(c[1]) > 0 {
				info.Other[c[0]] = c[1]
			}
		}
	}

	return Info{}
}
