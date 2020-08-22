package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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
	Run:     generateFunc,
}

func init() {
	generateCmd.Flags().StringVarP(&fileName, "name", "n", "grom.json", "the name of the generated grom configuration file")
	rootCmd.AddCommand(generateCmd)
}

func generateFunc(cmd *cobra.Command, args []string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(fileInfo))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fileInfo)
	fmt.Println("\nwrite in:", fileName)
}

func generateFileInfo() string {
	c := util.CMDConfig{
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
		EnableFieldComment: true,
		EnableSqlNull:      false,
		EnableGureguNull:   false,
		EnableJsonTag:      true,
		EnableXmlTag:       false,
		EnableGormTag:      true,
		EnableXormTag:      false,
		EnableBeegoTag:     false,
		EnableGoroseTag:    false,
	}

	b, _ := json.MarshalIndent(&c, "", "    ")
	fileInfo = string(b)

	return fileInfo
}
