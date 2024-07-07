// v0.2.0
// Author: DIEHL E.

package audio

import (
	"encoding/json"
	"io"
	"math/rand/v2"
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

// Len returns the number of tracks in the Collection.
func (c Collection) Len() int {
	return len(c.tracks)
}

func (c *Collection) Load(rd io.Reader) error {
	data, err := io.ReadAll(rd)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.tracks)
}

// Populate scans `path` and adds all the audio tracks to the collection.
func (c *Collection) Populate(path string) error {
	return filepath.WalkDir(path, func(path1 string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !isAudioFile(path1) {
			return nil
		}
		tr, err := NewTrack(path1)
		if err != nil {
			return nil
		}
		c.addTrack(*tr)
		return nil
	})
}

func (c Collection) Random() Track {
	i := rand.IntN(len(c.tracks))
	return c.tracks[i]
}

func (c Collection) Store(wr io.Writer) error {
	data, err := json.Marshal(c.tracks)
	if err != nil {
		return err
	}
	_, err = wr.Write(data)
	return err
}

// --------------------------

func (c *Collection) addTrack(t Track) {
	c.tracks = append(c.tracks, t)
}

func isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp3" || ext == ".wav" || ext == ".flac"
}
