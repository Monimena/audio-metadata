package parser

import (
	"github.com/Monimena/audio-metadata"
	"github.com/jfreymuth/oggvorbis"
	"github.com/jfreymuth/vorbis"
	"io"
	"strconv"
	"strings"
)

type VorbisParser struct {

}

func (p VorbisParser) Parse(file io.Reader) (*metadata.Info, error) {
	comments, err := oggvorbis.GetCommentHeader(file)
	info := mapMetadataVorbis(comments)

	if err != nil {
		return nil, err
	}

	return &info, nil
}

func mapMetadataVorbis(header vorbis.CommentHeader) metadata.Info {
	info := metadata.Info{}

	for _, comment := range header.Comments {
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
			info.Other[c[0]] = c[1]
		}
	}

	return metadata.Info{}
}