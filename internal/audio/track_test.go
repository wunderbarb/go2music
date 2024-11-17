// v0.2.0
// Author: DIEHL E.

package audio

import (
	"github.com/wunderbarb/test"
	"path/filepath"
	"testing"
)

func TestNewTrack(t *testing.T) {
	require, assert := test.Describe(t)
	const testDir = "testdata"

	tests := []struct {
		filePath   string
		expSuccess bool
	}{
		{"01. Autumn Rain.flac", true},
		{"dumb.flac", false},
		{"dumb.mp4", false},
		{"bad", false},
	}

	for _, tt := range tests {
		track, err := NewTrack(filepath.Join(testDir, tt.filePath))
		require.Equal(tt.expSuccess, err == nil)
		if err == nil {
			assert.NotNil(track)
		}
	}
}

func Test_extractDataFLAC(t *testing.T) {
	require, assert := test.Describe(t)
	album, title, err := extractDataFLAC(_goodTrack)
	require.NoError(err)
	assert.Equal("Blue Moon", album)
	assert.Equal("Autumn Rain", title)
}

func isPanic(err error) {
	if err != nil {
		panic(err)
	}
}
