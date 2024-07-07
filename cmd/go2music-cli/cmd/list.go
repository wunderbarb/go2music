// v0.1.0

package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/wunderbarb/go2music/internal/audio"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all devices",
	Long:  `list all devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getPlayer()
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(4)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func getPlayer() (*audio.Player, error) {
	var p audio.Player
	ll, err := p.Devices()
	if err != nil {
		return nil, err
	}
	answer, err := pterm.DefaultInteractiveSelect.WithOptions(ll).Show()
	if err != nil {
		return nil, err
	}
	if err := p.SelectDevice(answer); err != nil {
		return nil, err
	}
	return &p, nil
}
