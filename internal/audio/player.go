// v0.1.0

package audio

import (
	"errors"
	"github.com/alexballas/go2tv/devices"
)

var ErrNoDeviceAvailable = errors.New("no device available")

type Player struct {
	devices        map[string]string
	selectedDevice string
}

// Devices list
func (p *Player) Devices() ([]string, error) {
	err := p.available()
	if err != nil {
		return nil, err
	}
	if len(p.devices) == 0 {
		return nil, ErrNoDeviceAvailable
	}
	s := make([]string, 0, len(p.devices))
	for k, _ := range p.devices {
		s = append(s, k)
	}
	return s, nil
}

func (p *Player) SelectDevice(d string) error {
	_, ok := p.devices[d]
	if !ok {
		return ErrNoDeviceAvailable
	}
	p.selectedDevice = d
	return nil
}

func (p *Player) available() error {
	deviceList, err := devices.LoadSSDPservices(1)
	if err != nil {
		return err
	}
	p.devices = deviceList
	return nil
}
