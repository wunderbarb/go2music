package main

import (
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pterm/pterm"
	"github.com/wunderbarb/go2music/internal/audio"
	"os"
	"path/filepath"
)

//go:embed version.txt
var cVersion string

const _cDefaultDB = "music.db"

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("go2music " + cVersion)
	myWindow.Resize(fyne.NewSize(800, 500))
	var c audio.Collection
	hDir, err := os.Getwd()
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(5)
	}
	rd, err := os.Open(filepath.Join(hDir, _cDefaultDB))
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(5)
	}
	if err = c.Load(rd); err != nil {
		pterm.Error.Println(err)
		os.Exit(6)
	}
	var p audio.Player
	ll, err := p.Devices()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer p.TearDown()
	title := binding.NewString()
	label := widget.NewLabel("Hello World")
	label.Bind(title)
	usf := newUserFeedBack(c, &p, title)
	// Create a new set of radio buttons
	radio := widget.NewRadioGroup(ll, func(s string) {
		if err := p.SelectDevice(s); err != nil {
			pterm.Error.Println(err)
			os.Exit(7)
		}
		fmt.Println("Radio clicked " + s)
	})

	// Create buttons with icons
	var started bool
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		fmt.Println("Play clicked")
		if !started {
			tr := c.Random()
			_ = usf.title.Set(tr.Title + " of " + tr.Album)
			fmt.Println(usf.title)
			if err := p.PlayTrack(tr, &usf); err != nil {
				fmt.Println(err)
				os.Exit(8)
			}
			started = true
			return
		}
		if err := p.Play(); err != nil {
			pterm.Error.Println(err)
			os.Exit(7)
		}

	})
	pauseButton := widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func() {
		if err := p.Pause(); err != nil {
			pterm.Error.Println(err)
			os.Exit(7)
		}
		fmt.Println("Pause clicked")
	})
	nextButton := widget.NewButtonWithIcon("Next", theme.MediaSkipNextIcon(), func() {
		fmt.Println("Next clicked")
		usf.Fini()
	})

	// Create a container for the buttons
	buttonBar := container.NewHBox(playButton, pauseButton, nextButton)

	// Set the content of the window to the radio buttons and buttons bar
	myWindow.SetContent(container.NewBorder(nil, buttonBar, radio, label))

	// Show and run the application
	myWindow.ShowAndRun()
}

type userFeedBack struct {
	c     *audio.Collection
	pl    *audio.Player
	title binding.String
}

func newUserFeedBack(c audio.Collection, pl *audio.Player, t binding.String) userFeedBack {
	return userFeedBack{c: &c, pl: pl, title: t}
}

func (ufb *userFeedBack) EmitMsg(s string) {
	//const (
	//	playSymbol  = "▶️"
	//	pauseSymbol = "⏸️"
	//	eraseEOL    = "\033[K"
	//)
	switch s {
	case "Playing":
		fmt.Println("Playing")
		//_= ufb.title.Set(playSymbol + "  " + ufb.title + eraseEOL)
	case "Paused":
		//	_=ufb.title.Set(pauseSymbol + "  " + ufb.title + eraseEOL)
		fmt.Println("Paused")
	}
}

func (ufb *userFeedBack) Fini() {
	tr := ufb.c.Random()
	_ = ufb.pl.Next(tr)
	_ = ufb.title.Set(tr.Title + " of " + tr.Album)
}
