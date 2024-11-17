// v0.1.1
// (C) Nov 2024

package cmd

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

// embed:version.txt
var cVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cmd",
	Version: cVersion,
	Short:   "go2music is a specialized UPNP player",
	Long:    `go2music is a specialized UPNP player.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
