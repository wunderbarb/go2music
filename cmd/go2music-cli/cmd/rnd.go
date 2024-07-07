// v0.1.0

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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		keyboard.Listen(func(key keys.Key) (stop bool, err error) {
			switch key.Code {
			case keys.CtrlC:
				return true, nil // Return true to stop listener
			case keys.Space:
				pl.Pause()
				return false, nil
			case keys.RuneKey: // Check if key is a rune key (a, b, c, 1, 2, 3, ...)
				switch key.String() {
				case "q":
					return true, nil
				case "s":
					pl.Stop()
					return false, nil
				case "p":
					pl.Play()
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
