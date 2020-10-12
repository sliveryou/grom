package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	projectName    = "Grom"
	projectVersion = "1.0.1"
	goVersion      = "go1.14.3"
	gitCommit      = "9f35c08031"
	buildTime      = "2020-10-12 10:11:49"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the grom version information",
	Long:  "Show the grom version information, such as project name, project version, go version, git commit id, build time, etc",
	Run:   versionFunc,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("Cmd Tool: %s\n", projectName)
	fmt.Printf(" Version: %s\n", projectVersion)
	fmt.Printf(" Go Version: %s\n", goVersion)
	fmt.Printf(" Git Commit: %s\n", gitCommit)
	fmt.Printf(" Build Time: %s\n", buildTime)
}
