package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/sliveryou/grom/util"
)

var (
	fileInfo string
	fileName string
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate grom configuration file",
	Long:    fmt.Sprintf("Generate grom configuration file like this:\n%s", generateFileInfo()),
	Example: "  grom generate -n ./grom.json",
	RunE:    generateFunc,
}

func init() {
	generateCmd.Flags().StringVarP(&fileName, "name", "n", "grom.json", "the name of the generated grom configuration file")
	rootCmd.AddCommand(generateCmd)
}

func generateFunc(_ *cobra.Command, _ []string) error {
	err := os.WriteFile(fileName, []byte(fileInfo), writeFilePerm)
	if err != nil {
		return errors.WithMessage(err, "os.WriteFile err")
	}

	fmt.Println("write in:", fileName)

	return nil
}

func generateFileInfo() string {
	c := util.CmdConfig{
		DBConfig: util.DBConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "user",
			Password: "password",
			Database: "database",
			Table:    "table",
		},
		PackageName:        "package_name",
		StructName:         "struct_name",
		EnableInitialism:   true,
		EnableFieldComment: true,
		EnableSQLNull:      false,
		EnableGureguNull:   false,
		EnableJSONTag:      true,
		EnableXMLTag:       false,
		EnableGormTag:      false,
		EnableXormTag:      false,
		EnableBeegoTag:     false,
		EnableGoroseTag:    false,
		EnableGormV2Tag:    true,
		DisableUnsigned:    false,
	}

	b, _ := json.MarshalIndent(&c, "", "    ")
	fileInfo = string(b)

	return fileInfo
}
