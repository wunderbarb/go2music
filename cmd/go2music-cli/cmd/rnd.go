// v0.1.0

package cmd

import (
	"context"
	"github.com/pterm/pterm"
	"github.com/wunderbarb/go2music/internal/audio"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
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
		exitCTX, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		scr := audio.NewDummyScreen(cancel)
		tr := c.Random()
		if err := pl.PlayTrack(tr, scr); err != nil {
			pterm.Error.Println(err)
			os.Exit(8)
		}
		pterm.Info.Println(tr.FilePath)
		<-exitCTX.Done()
		pl.TearDown()
	},
}

func init() {
	rootCmd.AddCommand(rndCmd)

}
