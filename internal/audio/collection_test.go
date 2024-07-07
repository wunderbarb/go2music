// v0.3.0
// Author: DIEHL E.
// (C), Jul 2024

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

func TestCollection_Load(t *testing.T) {
	require, assert := test.Describe(t)

	c := NewCollection()
	_ = c.Populate("testdata")
	wr := test.NewInRAMWriter()

	require.NoError(c.Store(wr))
	require.NoError(wr.Close())

	c1 := NewCollection()
	rd := test.NewInRAMReader(wr.Bytes())
	require.NoError(c1.Load(rd))
	assert.Equal(c1.Len(), c.Len())
}
