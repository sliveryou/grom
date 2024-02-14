package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	projectName    = "grom"
	projectVersion = "1.0.7"
	goVersion      = ""
	gitCommit      = ""
	buildTime      = ""
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show the grom version information",
	Long:    "Show the grom version information, such as project name, project version, go version, git commit id, build time, etc",
	Example: "  grom version",
	RunE:    versionFunc,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionFunc(_ *cobra.Command, _ []string) error {
	fmt.Printf("%s:\n", projectName)
	fmt.Printf("   version: %s %s/%s\n", projectVersion, runtime.GOOS, runtime.GOARCH)
	if goVersion == "" {
		goVersion = runtime.Version()
	}
	fmt.Printf("   go version: %s\n", goVersion)
	if gitCommit != "" {
		fmt.Printf("   git commit: %s\n", gitCommit)
	}
	if buildTime != "" {
		fmt.Printf("   build time: %s\n", buildTime)
	}

	return nil
}
