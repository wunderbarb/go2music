// v0.1.0

package audio

import (
	"path/filepath"
	"testing"

	"github.com/wunderbarb/test"
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
