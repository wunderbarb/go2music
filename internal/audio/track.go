// v0.1.1
// Author: DIEHL E.

package audio

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alexballas/go2tv/soapcalls/utils"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/pkg/errors"
)

// ErrNotAudioFile is returned when the file is not an audio file.
var ErrNotAudioFile = errors.New("not an audio file")

var ErrNoInformation = errors.New("no information")

// Track holds the metadata for a music track.
type Track struct {
	FilePath  string `json:"file_path"`
	MediaType string `json:"media_type"`
	// Album to which the track belongs
	Album string `json:"album"`
	Title string `json:"title"`
}

// NewTrack creates a new Track.  It certifies that the file is an audio file.
func NewTrack(filePath string) (*Track, error) {
	if !isAudioFile(filePath) {
		return nil, ErrNotAudioFile
	}
	absMediaFile, err := filepath.Abs(filePath)
	if err != nil {
		return nil, errors.WithMessage(err, "absolute path")
	}
	mFile, err := os.Open(absMediaFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "open file")
	}
	defer func() { _ = mFile.Close() }()
	mediaType, err := utils.GetMimeDetailsFromFile(mFile)
	if err != nil {
		return nil, errors.WithMessage(err, "get mime details")
	}
	if !strings.HasPrefix(mediaType, "audio/") {
		return nil, ErrNotAudioFile
	}
	tr := &Track{FilePath: absMediaFile, MediaType: mediaType}
	switch filepath.Ext(tr.FilePath) {
	case ".flac":
		a, t, err := extractDataFLAC(filePath)
		if err == nil {
			tr.Title = t
			tr.Album = a
		}
	}
	return tr, nil
}

func extractDataFLAC(path string) (string, string, error) {
	const (
		// https://www.xiph.org/vorbis/doc/v-comment.html
		cTitle = "TITLE"
		cAlbum = "ALBUM"
	)
	rd, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	bl, err := flac.ParseMetadata(rd)
	if err != nil {
		return "", "", err
	}
	for _, block := range bl.Meta {
		if block.Type != flac.VorbisComment {
			continue
		}
		vorbisComment, err := flacvorbis.ParseFromMetaDataBlock(*block)
		if err != nil {
			return "", "", err
		}
		s, err := vorbisComment.Get(cTitle)
		var title, album string
		if err == nil && len(s) > 0 {
			title = s[0]
		}
		s, err = vorbisComment.Get(cAlbum)
		if err == nil && len(s) > 0 {
			album = s[0]
		}
		return album, title, err
	}
	return "", "", ErrNoInformation
}
