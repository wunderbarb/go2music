// v0.2.0
// Author: DIEHL E.

package audio

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexballas/go2tv/devices"
	"github.com/alexballas/go2tv/httphandlers"
	"github.com/alexballas/go2tv/soapcalls"
)

var (
	ErrNoDeviceAvailable = errors.New("no device available")
	ErrNoDevicePlaying   = errors.New("no device playing")
)

type Player struct {
	devices        map[string]string
	selectedDevice string
	server         *httphandlers.HTTPserver
	tvData         *soapcalls.TVPayload
}

// Devices lists the UPNP renderer currently available.
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

// PlayTrack plays the track `tr` on the selected device.  `src` handles the feedback and stopping the service.
func (p *Player) PlayTrack(tr Track, scr httphandlers.Screen) error {
	tvData, err := soapcalls.NewTVPayload(&soapcalls.Options{
		DMR:       p.devices[p.selectedDevice],
		Media:     tr.FilePath,
		Subs:      "",
		Mtype:     tr.MediaType,
		Transcode: false,
		Seek:      false,
		LogOutput: nil,
	})
	if err != nil {
		return err
	}
	p.tvData = tvData
	p.server = httphandlers.NewServer(p.tvData.ListenAddress())
	serverStarted := make(chan error)
	// We pass the tvData here as we need the callback handlers to be able to react
	// to the different media renderer states.
	go func() {
		p.server.StartServer(serverStarted, tr.FilePath, "", p.tvData, scr)
	}()
	// Wait for HTTP server to properly initialize
	if err := <-serverStarted; err != nil {
		return err
	}
	return p.tvData.SendtoTV("Play1")
}

// SelectDevice defines the renderer device to be used amongst the available devices.
func (p *Player) SelectDevice(d string) error {
	_, ok := p.devices[d]
	if !ok {
		return ErrNoDeviceAvailable
	}
	p.selectedDevice = d
	return nil
}

func (p *Player) Stop() error {
	if p.tvData == nil {
		return ErrNoDevicePlaying
	}
	return p.tvData.SendtoTV("Stop")
}

// ------------------------------------------

func (p *Player) available() error {
	deviceList, err := devices.LoadSSDPservices(1)
	if err != nil {
		return err
	}
	p.devices = deviceList
	return nil
}

func (p *Player) TearDown() {
	_ = p.Stop()
	p.tvData = nil
	if p.server != nil {
		p.server.StopServer()
	}
}

type DummyScreen struct {
	ctxCancel context.CancelFunc
}

func NewDummyScreen(cancel context.CancelFunc) *DummyScreen {
	return &DummyScreen{ctxCancel: cancel}
}

func (d *DummyScreen) EmitMsg(msg string) {
	fmt.Println(msg)
}

func (d *DummyScreen) Fini() {
	d.ctxCancel()
}
