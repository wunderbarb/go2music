// v0.2.0
// Author: DIEHL E.

package audio

import (
	"github.com/wunderbarb/test"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const _cLocalRenderer = "Audiophile UPnP Renderer"

var (
	_goodTrack1 = filepath.Join("testdata", "album", "02. Blue Moon.flac")
	_goodTrack  = filepath.Join("testdata", "01. Autumn Rain.flac")
)

func TestPlayer_available(t *testing.T) {
	_, assert := test.Describe(t)

	var p Player
	assert.NoError(p.available())

}

func TestPlayer_Devices(t *testing.T) {
	require, assert := test.Describe(t)

	var p Player
	devices, err := p.Devices()
	require.NoError(err)
	assert.NotZero(len(devices))
}

func TestPlayer_SelectDevice(t *testing.T) {
	_, assert := test.Describe(t)

	var p Player
	devices, _ := p.Devices()

	assert.NoError(p.SelectDevice(devices[0]))
	assert.Error(p.SelectDevice("bad"))

}

func TestPlayer_PlayTrack(t *testing.T) {
	require, _ := test.Describe(t)
	defer test.NoLeakButPersistentHTTP(t)

	pt := getTestPlayer()
	src := &DummyScreen{}
	track, _ := NewTrack(_goodTrack)
	defer pt.TearDown()

	require.NoError(pt.PlayTrack(*track, src))
	time.Sleep(5 * time.Second)
}

func TestPlayer_Next(t *testing.T) {
	require, _ := test.Describe(t)
	defer test.NoLeakButPersistentHTTP(t)

	pt := getTestPlayer()
	src := &DummyScreen{}
	track, _ := NewTrack(_goodTrack)
	defer pt.TearDown()
	track1, _ := NewTrack(_goodTrack1)

	require.NoError(pt.PlayTrack(*track, src))
	time.Sleep(5 * time.Second)
	require.NoError(pt.Next(*track1))
	time.Sleep(5 * time.Second)
}

// -----------------------------------------

func getTestPlayer() Player {
	var p Player
	devices, _ := p.Devices()
	for _, device := range devices {
		if strings.HasPrefix(device, _cLocalRenderer) {
			_ = p.SelectDevice(device)
			break
		}
	}
	return p
}
