package cmd

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sliveryou/grom/pkg/api"
)

var configPath string

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate go-zero api project by mysql information schema",
	Long:    "Generate go-zero api project by mysql information_schema.columns and information_schema.statistics",
	Example: "  grom api generate -n ./config.yaml",
	RunE:    generateFunc,
}

func init() {
	generateCmd.Flags().StringVarP(&configPath, "name", "n", "config.yaml", "the name of the grom api configuration file")
	APICmd.AddCommand(generateCmd)
}

func generateFunc(_ *cobra.Command, _ []string) error {
	config, err := parseConfig()
	if err != nil {
		return errors.WithMessage(err, "parseConfig err")
	}

	return errors.WithMessage(api.GenerateProject(config), "api.GenerateProject err")
}

func parseConfig() (api.ProjectConfig, error) {
	config := api.ProjectConfig{}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return api.ProjectConfig{}, errors.WithMessage(err, "viper.ReadInConfig err")
	}

	if err := viper.Unmarshal(&config, func(c *mapstructure.DecoderConfig) {
		c.Squash = true
		c.TagName = "json"
	}); err != nil {
		return api.ProjectConfig{}, errors.WithMessage(err, "viper.Unmarshal err")
	}

	return config, nil
}
