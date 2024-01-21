package cmd

import (
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

const (
	// codeFailure failure code
	codeFailure = 1
)

var rootCmd = &cobra.Command{
	Use:   "grom",
	Short: "Get golang model structure by mysql information schema",
	Example: "  grom generate -n ./grom.json\n" +
		"  grom convert -n ./grom.json",
}

// Execute executes the root command and its subcommands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red.Println(err.Error())
		os.Exit(codeFailure)
	}
}
