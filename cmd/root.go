package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		fmt.Println(err)
		os.Exit(1)
	}
}
