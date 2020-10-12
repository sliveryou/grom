package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sliveryou/grom/util"
)

var (
	filePath    string
	packageName string
	structName  string
	host        string
	port        int
	user        string
	password    string
	database    string
	table       string
	enable      []string

	validServices = map[string]struct{}{
		"INITIALISM":    {},
		"FIELD_COMMENT": {},
		"SQL_NULL":      {},
		"GUREGU_NULL":   {},
		"JSON_TAG":      {},
		"XML_TAG":       {},
		"GORM_TAG":      {},
		"XORM_TAG":      {},
		"BEEGO_TAG":     {},
		"GOROSE_TAG":    {},
		"GORM_V2_TAG":   {},
	}
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert mysql table fields to golang model structure",
	Long:  "Convert mysql table fields to golang model structure by information_schema.columns and information_schema.statistics",
	Example: "  grom convert -n ./grom.json\n" +
		"  grom convert -H localhost -P 3306 -u user -p password -d database -t table -e INITIALISM,FIELD_COMMENT,JSON_TAG,GORM_V2_TAG --package PACKAGE_NAME --struct STRUCT_NAME",
	Run: convertFunc,
}

func init() {
	convertCmd.Flags().StringVar(&packageName, "package", "", "the package name of the converted model structure")
	convertCmd.Flags().StringVar(&structName, "struct", "", "the struct name of the converted model structure")
	convertCmd.Flags().StringVarP(&filePath, "name", "n", "", "the name of the grom configuration file")
	convertCmd.Flags().StringVarP(&host, "host", "H", "", "the host of mysql")
	convertCmd.Flags().IntVarP(&port, "port", "P", 0, "the port of mysql")
	convertCmd.Flags().StringVarP(&user, "user", "u", "", "the user of mysql")
	convertCmd.Flags().StringVarP(&password, "password", "p", "", "the password of mysql")
	convertCmd.Flags().StringVarP(&database, "database", "d", "", "the database of mysql")
	convertCmd.Flags().StringVarP(&table, "table", "t", "", "the table of mysql")
	convertCmd.Flags().StringSliceVarP(&enable, "enable", "e", nil, "enable services (must in [INITIALISM,FIELD_COMMENT,SQL_NULL,GUREGU_NULL,JSON_TAG,XML_TAG,GORM_TAG,XORM_TAG,BEEGO_TAG,GOROSE_TAG,GORM_V2_TAG])")

	rootCmd.AddCommand(convertCmd)
}

func convertFunc(cmd *cobra.Command, args []string) {
	config, err := getCmdConfig()
	if err != nil {
		return
	}

	s, err := util.ConvertTable(config)
	if err != nil {
		return
	}

	fmt.Println(s)
}

func getCmdConfig() (util.CMDConfig, error) {
	config := util.CMDConfig{}

	if filePath != "" {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err)
			return util.CMDConfig{}, err
		}

		err = json.Unmarshal(content, &config)
		if err != nil {
			fmt.Println(err)
			return util.CMDConfig{}, err
		}
	}

	if packageName != "" {
		config.PackageName = packageName
	}
	if structName != "" {
		config.StructName = structName
	}
	if host != "" {
		config.Host = host
	}
	if port != 0 {
		config.Port = port
	}
	if user != "" {
		config.User = user
	}
	if password != "" {
		config.Password = password
	}
	if database != "" {
		config.Database = database
	}
	if table != "" {
		config.Table = table
	}

	if len(enable) != 0 {
		for _, v := range enable {
			service := strings.ToUpper(v)
			if _, ok := validServices[service]; !ok {
				err := errors.New("enabled service is invalid, service: " + service)
				fmt.Println(err)
				return util.CMDConfig{}, err
			}

			switch service {
			case "INITIALISM":
				config.EnableInitialism = true
			case "FIELD_COMMENT":
				config.EnableFieldComment = true
			case "SQL_NULL":
				config.EnableSqlNull = true
			case "GUREGU_NULL":
				config.EnableGureguNull = true
			case "JSON_TAG":
				config.EnableJsonTag = true
			case "XML_TAG":
				config.EnableXmlTag = true
			case "GORM_TAG":
				config.EnableGormTag = true
			case "XORM_TAG":
				config.EnableXormTag = true
			case "BEEGO_TAG":
				config.EnableBeegoTag = true
			case "GOROSE_TAG":
				config.EnableGoroseTag = true
			case "GORM_V2_TAG":
				config.EnableGormV2Tag = true
			}
		}
	}

	return config, nil
}
