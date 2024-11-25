package main

import (
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
	img := canvas.NewImageFromFile("go2music logo.png")
	img.Resize(fyne.NewSize(400, 400))
	img.FillMode = canvas.ImageFillContain
	usf := newUserFeedBack(c, &p, title, img)
	// Create a new set of radio buttons
	radio := widget.NewRadioGroup(ll, func(s string) {
		if err := p.SelectDevice(s); err != nil {
			pterm.Error.Println(err)
			os.Exit(7)
		}
	})

	// Create buttons with icons
	var started bool
	playButton := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		fmt.Println("Play clicked")
		if !started {
			tr := c.Random()
			usf.update(&tr)
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
	})
	nextButton := widget.NewButtonWithIcon("Next", theme.MediaSkipNextIcon(), func() {
		usf.Fini()
	})

	// Create a container for the buttons
	buttonBar := container.NewHBox(playButton, pauseButton, nextButton)

	// Set the content of the window to the radio buttons and buttons bar
	myWindow.SetContent(container.NewBorder(label, buttonBar, radio, nil, img))

	// Show and run the application
	myWindow.ShowAndRun()
}

type userFeedBack struct {
	c     *audio.Collection
	pl    *audio.Player
	title binding.String
	img   *canvas.Image
}

func newUserFeedBack(c audio.Collection, pl *audio.Player, t binding.String, img *canvas.Image) userFeedBack {
	return userFeedBack{c: &c, pl: pl, title: t, img: img}
}

func (ufb *userFeedBack) EmitMsg(s string) {
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
	ufb.update(&tr)
}

func (ufb userFeedBack) update(tr *audio.Track) {
	_ = ufb.title.Set(tr.Title + " of " + tr.Album)
	cover, err := tr.Cover()
	if err != nil {
		cover = "go2music logo.png"
	}
	ufb.img.File = cover
	ufb.img.Refresh()
}
