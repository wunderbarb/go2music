// v0.1.0

package audio

import (
	"github.com/wunderbarb/test"
	"testing"
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
