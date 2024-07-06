// v0.2.0
// Author: DIEHL E.

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

func TestCollection_Populate(t *testing.T) {
	require, assert := test.Describe(t)

	c := NewCollection()
	require.NoError(c.Populate("testdata"))
	assert.Equal(2, c.Len())

}
