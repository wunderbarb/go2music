package audio

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/alexballas/go2tv/soapcalls/utils"
)

// ErrNotAudioFile is returned when the file is not an audio file.
var ErrNotAudioFile = errors.New("not an audio file")

// Track holds the metadata for a music track.
type Track struct {
	filePath  string
	mediaType string
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
	return &Track{filePath: absMediaFile, mediaType: mediaType}, nil
}
