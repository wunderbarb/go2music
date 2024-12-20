// v0.2.1
// Author: DIEHL E.
// (C), Jul 2024

package audio

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexballas/go2tv/devices"
	"github.com/alexballas/go2tv/httphandlers"
	"github.com/alexballas/go2tv/soapcalls"
	"github.com/alexballas/go2tv/soapcalls/utils"
	"github.com/pterm/pterm"
	"net/url"
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
	for k := range p.devices {
		s = append(s, k)
	}
	return s, nil
}

// PlayTrack plays the track `tra` on the selected device.  `src` handles the feedback and stopping the service.
func (p *Player) PlayTrack(tra Track, scr httphandlers.Screen) error {
	tvData, err := soapcalls.NewTVPayload(&soapcalls.Options{
		DMR:       p.devices[p.selectedDevice],
		Media:     tra.FilePath,
		Subs:      "",
		Mtype:     tra.MediaType,
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
		p.server.StartServer(serverStarted, tra.FilePath, "", p.tvData, scr)
	}()
	// Wait for HTTP server to properly initialize
	if err := <-serverStarted; err != nil {
		return err
	}
	return p.tvData.SendtoTV("Play1")
}

func (p *Player) Pause() error {
	if p.tvData == nil {
		return ErrNoDevicePlaying
	}
	return p.tvData.SendtoTV("Pause")
}

func (p *Player) Play() error {
	if p.tvData == nil {
		return ErrNoDevicePlaying
	}
	return p.tvData.SendtoTV("Play1")
}

// Next plays the next track `tra`
func (p *Player) Next(tra Track) error {
	if p.tvData == nil {
		return ErrNoDevicePlaying
	}
	oldHandler, err := url.Parse(p.tvData.MediaURL)
	if err != nil { // SHOULD NEVER HAPPEN
		return err
	}
	listenAddress := p.tvData.ListenAddress()
	p.tvData.MediaURL = fmt.Sprintf("http://%s/%s", listenAddress, utils.ConvertFilename(tra.FilePath))
	p.tvData.MediaPath = tra.FilePath
	p.tvData.MediaType = tra.MediaType
	newHandler, err := url.Parse(p.tvData.MediaURL)
	if err != nil { // SHOULD NEVER HAPPEN
		return err
	}
	p.server.AddHandler(newHandler.Path, p.tvData, tra.FilePath)
	p.server.RemoveHandler(oldHandler.Path)
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
	pterm.Info.Println(msg)
}

func (d *DummyScreen) Fini() {
	d.ctxCancel()
}
