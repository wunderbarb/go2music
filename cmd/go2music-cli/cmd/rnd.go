// v0.1.1
// Author: DIEHL E.

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

		scr := audio.NewDummyScreen(nil)
		tr := c.Random()
		if err := pl.PlayTrack(tr, scr); err != nil {
			pterm.Error.Println(err)
			os.Exit(8)
		}
		pterm.Info.Println(tr.FilePath)
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
					tr := c.Random()
					_ = pl.Next(tr)
					pterm.Info.Println(tr.FilePath)
					return false, nil
				}
			}

			return false, nil // Return false to continue listening
		})
		pl.TearDown()
	},
}

func init() {
	rootCmd.AddCommand(rndCmd)

}
