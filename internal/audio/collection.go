// v0.1.0

package audio

import (
	"os"
	"path/filepath"
	"strings"
)

// Collection holds all the music.
type Collection struct {
	tracks []Track
}

func NewCollection() *Collection {
	return &Collection{}
}

func (c *Collection) Populate(path string) error {
	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !isAudioFile(path) {
			return nil
		}
		c.addTrack(Track{filePath: path})
		return nil
	})
}

func (c *Collection) addTrack(t Track) {
	c.tracks = append(c.tracks, t)
}

func isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp3" || ext == ".wav" || ext == ".flac"
}
