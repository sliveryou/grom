package cmd

import "github.com/spf13/cobra"

const (
	// writeFilePerm default perm for writing file.
	writeFilePerm = 0o666
)

// APICmd represents the api root command.
var APICmd = &cobra.Command{
	Use:   "api",
	Short: "Get go-zero api project by mysql information schema",
	Example: "  grom api config -n ./api-config.json\n" +
		"  grom api generate -n ./api-config.json",
}
