package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	configName string
	//go:embed config.yaml
	configInfoBytes []byte
)

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Generate grom api configuration file",
	Long:    "Generate grom api configuration file like this:\n" + string(configInfoBytes),
	Example: "  grom api config -n ./config.yaml",
	RunE:    configFunc,
}

func init() {
	configCmd.Flags().StringVarP(&configName, "name", "n", "config.yaml", "the name of the generated grom api configuration file")
	APICmd.AddCommand(configCmd)
}

func configFunc(_ *cobra.Command, _ []string) error {
	err := os.WriteFile(configName, configInfoBytes, writeFilePerm)
	if err != nil {
		return errors.WithMessage(err, "os.WriteFile err")
	}

	fmt.Println("write in:", configName)

	return nil
}
