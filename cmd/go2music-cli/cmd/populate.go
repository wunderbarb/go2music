// v0.1.0

package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"

	"github.com/wunderbarb/go2music/internal/audio"
)

var output string

const _cDefaultDB = "music.db"

// populateCmd represents the populate command
var populateCmd = &cobra.Command{
	Use:   "populate",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := audio.NewCollection()
		if err := c.Populate(args[0]); err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
		wr, err := os.Create(output)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(2)
		}
		defer wr.Close()
		if err := c.Store(wr); err != nil {
			pterm.Error.Println(err)
			os.Exit(3)
		}
		pterm.Success.Printfln("Listed %d songs", c.Len())
	},
}

func init() {
	rootCmd.AddCommand(populateCmd)

	populateCmd.Flags().StringVarP(&output, "output", "o", _cDefaultDB, "Output filename")
}
