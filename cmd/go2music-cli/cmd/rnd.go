// v0.1.4
// Author: DIEHL E.
// (C) Nov 2024

package cmd

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/wunderbarb/go2music/internal/audio"
	"os"
)

// rndCmd represents the rnd command
var rndCmd = &cobra.Command{
	Use:   "rnd",
	Short: "plays randomly tracks from the collection.",
	Long:  `plays randomly tracks from the collection.`,
	Run: func(cmd *cobra.Command, args []string) {
		var c audio.Collection
		rd, err := os.Open(_cDefaultDB)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(5)
		}
		if err = c.Load(rd); err != nil {
			pterm.Error.Println(err)
			os.Exit(6)
		}
		pl, err := getPlayer()
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(7)
		}
		defer pl.TearDown()
		//scr := audio.NewDummyScreen(func() { pterm.Warning.Println("dummy") })
		scr := newUserFeedBack(c, pl)
		tr := c.Random()
		scr.title = tr.Title + " of " + tr.Album
		if err := pl.PlayTrack(tr, &scr); err != nil {
			pterm.Error.Println(err)
			os.Exit(8)
		}

		_ = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
			switch key.Code {
			case keys.CtrlC:
				return true, nil // Return true to stop listener
			case keys.Space:
				_ = pl.Pause()
				return false, nil
			case keys.RuneKey: // Check if key is a rune key (a, b, c, 1, 2, 3, ...)
				switch key.String() {
				case "q":
					return true, nil
				case "s":
					_ = pl.Stop()
					return false, nil
				case "p":
					_ = pl.Play()
					return false, nil
				case "n":
					scr.Fini()
					return false, nil
				}
			}

			return false, nil // Return false to continue listening
		})

	},
}

func init() {
	rootCmd.AddCommand(rndCmd)

}

type userFeedBack struct {
	c     *audio.Collection
	pl    *audio.Player
	sp    *pterm.SpinnerPrinter
	title string
}

func newUserFeedBack(c audio.Collection, pl *audio.Player) userFeedBack {
	s, _ := pterm.DefaultSpinner.Start()
	s.ShowTimer = false
	return userFeedBack{c: &c, pl: pl, sp: s}
}

func (ufb *userFeedBack) EmitMsg(s string) {
	const (
		playSymbol  = "▶️"
		pauseSymbol = "⏸️"
		eraseEOL    = "\033[K"
	)
	switch s {
	case "Playing":
		ufb.sp.UpdateText(playSymbol + "  " + ufb.title + eraseEOL)
	case "Paused":

		ufb.sp.UpdateText(pauseSymbol + "  " + ufb.title + eraseEOL)

	}
}

func (ufb *userFeedBack) Fini() {
	tr := ufb.c.Random()
	_ = ufb.pl.Next(tr)
	ufb.title = tr.Title + " of " + tr.Album
}
