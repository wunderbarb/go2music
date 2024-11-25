// v0.3.0
// Author: DIEHL E.
// (C), Nov 2024

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

func Test_isImageFile(t *testing.T) {
	_, assert := test.Describe(t)

	tests := []struct {
		name       string
		expSuccess bool
	}{
		{test.RandomID() + ".jpg", true},
		{test.RandomID() + ".png", true},
		{test.RandomID() + ".gif", true},
		{test.RandomID() + ".bmp", true},
		{test.RandomID() + ".bad", false},
	}
	for _, tt := range tests {
		assert.Equal(isImageFile(tt.name), tt.expSuccess)
	}
}

func TestTrack_Cover(t *testing.T) {
	require, assert := test.Describe(t)

	tests := []struct {
		path       string
		expSuccess bool
	}{
		{filepath.Join("testdata", "01. Autumn Rain.flac"), true},
		{filepath.Join("testdata", "album", "02. Blue Moon.flac"), true},
		{"track_test.go", false},
	}
	for i, tt := range tests {
		a, err := findCover(tt.path)
		require.Equal(tt.expSuccess, err == nil, "sample %d", i+1)
		if err == nil {
			assert.Contains(a, "go2music logo.png")
		}
	}
}
