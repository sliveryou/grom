package api

import (
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/sliveryou/grom/util"
)

// Config represents the config of the generated api.
type Config struct {
	util.DBConfig

	StructName       string   `json:"struct_name"`       // camel
	SnakeStructName  string   `json:"snake_struct_name"` // snake
	TableComment     string   `json:"-"`
	EnableInitialism bool     `json:"enable_initialism"`
	IgnoreFields     []string `json:"ignore_fields"`

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Author  string `json:"author"`
	Email   string `json:"email"`
	Version string `json:"version"`

	ServiceName  string `json:"service_name"`  // camel
	ServerPrefix string `json:"server_prefix"` // lower
	GroupPrefix  string `json:"group_prefix"`  // lower
}

// GetCMDConfig gets the *util.CMDConfig.
func (c *Config) GetCMDConfig() *util.CMDConfig {
	return &util.CMDConfig{
		DBConfig:         c.DBConfig,
		StructName:       c.StructName,
		EnableInitialism: c.EnableInitialism,
		DisableUnsigned:  true,
	}
}

// UpdateBy Config updated by util.CMDConfig.
func (c *Config) UpdateBy(cc *util.CMDConfig) {
	c.StructName = cc.StructName
	c.TableComment = cc.TableComment
}

// ProjectConfig represents the config of the generated project.
type ProjectConfig struct {
	Config
	Dir                 string   `json:"dir"`
	TablePrefix         string   `json:"table_prefix"`
	Tables              []string `json:"tables"`
	NeedTrimTablePrefix bool     `json:"need_trim_table_prefix"`
}

type generateConfig struct {
	IdComment       string
	SnakeStructName string // snake
	ModelName       string // camel
	GroupName       string // lower
	StructFields    []util.StructField
}

func getGenerateConfig(c *Config, fs []*util.StructField) generateConfig {
	gc := generateConfig{
		IdComment:       "id",
		SnakeStructName: strcase.ToSnake(c.StructName),
		ModelName:       c.StructName,
		GroupName:       strings.ToLower(c.StructName),
	}
	if c.SnakeStructName != "" {
		gc.SnakeStructName = c.SnakeStructName
		gc.ModelName = strcase.ToCamel(c.SnakeStructName)
		gc.GroupName = strings.ToLower(strcase.ToCamel(c.SnakeStructName))
	}

	fields := make([]util.StructField, 0, len(fs))
	for _, f := range fs {
		fi := *f
		// ignore fields
		if len(c.IgnoreFields) > 0 && contains(c.IgnoreFields, fi.RawName) {
			continue
		}
		// get id comment
		if fi.IsPrimaryKey {
			if fi.Comment != "" {
				gc.IdComment = fi.Comment
			} else {
				gc.IdComment = c.TableComment + "id"
			}
		}
		// remove unsigned
		if strings.HasPrefix(fi.Type, "u") {
			fi.Type = strings.TrimPrefix(fi.Type, "u")
		}
		// convert time.Time to int64
		if fi.Type == "time.Time" {
			fi.Type = "int64"
		}
		if enums := getEnums(fi.Comment); fi.Type == "bool" && enums == "0 1" {
			// trim bool comment
			fi.Comment = convertComment(fi.Comment, true)
		} else if fi.Type == "int" {
			if enums != "" {
				// convert enums int to int32
				fi.Type = "int32"
			} else {
				// convert int to int64 adapted to protobuf
				fi.Type = "int64"
			}
		}
		fields = append(fields, fi)
	}
	gc.StructFields = fields

	return gc
}
