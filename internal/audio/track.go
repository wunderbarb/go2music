// v0.1.0

package audio

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alexballas/go2tv/soapcalls/utils"
	"github.com/pkg/errors"
)

// ErrNotAudioFile is returned when the file is not an audio file.
var ErrNotAudioFile = errors.New("not an audio file")

// Track holds the metadata for a music track.
type Track struct {
	FilePath  string `json:"file_path"`
	MediaType string `json:"media_type"`
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
	return &Track{FilePath: absMediaFile, MediaType: mediaType}, nil
}
