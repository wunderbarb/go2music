// v0.1.0

package audio

import (
	"testing"

	"github.com/wunderbarb/test"
)

func Test_isAudioFile(t *testing.T) {
	_, assert := test.Describe(t)

	tests := []struct {
		path      string
		expResult bool
	}{
		{"file.mp3", true},
		{"file.wav", true},
		{"file.flac", true},
		{"file.mp4", false},
		{"file.Wav", true},
	}
	for _, tt := range tests {
		assert.Equal(isAudioFile(tt.path), tt.expResult)
	}
}
