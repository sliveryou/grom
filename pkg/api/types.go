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

// GetCmdConfig gets the *util.CmdConfig.
func (c *Config) GetCmdConfig() *util.CmdConfig {
	return &util.CmdConfig{
		DBConfig:         c.DBConfig,
		StructName:       c.StructName,
		EnableInitialism: c.EnableInitialism,
		DisableUnsigned:  true,
	}
}

// UpdateBy Config updated by util.CmdConfig.
func (c *Config) UpdateBy(cc *util.CmdConfig) {
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

// StructField represents the field of the generated model structure.
type StructField struct {
	Name         string
	Type         string
	Comment      string
	RawName      string
	RawType      string
	Default      string
	Enums        string
	IsPrimaryKey bool
	IsNullable   bool
}

// ToStructField converts the util.StructField to StructField.
func ToStructField(sf *util.StructField) StructField {
	return StructField{
		Name:         sf.Name,
		Type:         sf.Type,
		Comment:      sf.Comment,
		RawName:      sf.RawName,
		RawType:      sf.Type,
		Default:      sf.Default,
		Enums:        getEnums(sf.Comment),
		IsPrimaryKey: sf.IsPrimaryKey,
		IsNullable:   sf.IsNullable,
	}
}

// IsTimeField reports whether f is time field.
func IsTimeField(f StructField) bool {
	return f.RawType == util.GoTime &&
		(f.Type == util.GoInt64 || f.Type == util.GoTime || f.Type == util.GoPointerTime)
}

// IsAutoTimeField reports whether f is auto time field.
func IsAutoTimeField(f StructField) bool {
	return strings.Contains(f.RawName, autoTimeSuffix) &&
		f.Default == defaultCurrentTimestamp &&
		(f.Type == util.GoInt64 || f.Type == util.GoTime || f.Type == util.GoPointerTime)
}

type generateConfig struct {
	IdName          string
	IdType          string
	IdComment       string
	IdRawName       string
	SnakeStructName string // snake
	ModelName       string // camel
	GroupName       string // lower
	StructFields    []StructField
}

func getGenerateConfig(c *Config, fs []*util.StructField) generateConfig {
	gc := generateConfig{
		IdComment:       c.TableComment + defaultIdComment,
		SnakeStructName: strcase.ToSnake(c.StructName),
		ModelName:       c.StructName,
		GroupName:       strings.ToLower(c.StructName),
	}
	if c.SnakeStructName != "" {
		gc.SnakeStructName = c.SnakeStructName
		gc.ModelName = strcase.ToCamel(c.SnakeStructName)
		gc.GroupName = strings.ToLower(strcase.ToCamel(c.SnakeStructName))
	}

	fields := make([]StructField, 0, len(fs))
	for _, f := range fs {
		fi := ToStructField(f)
		// ignore fields
		if len(c.IgnoreFields) > 0 && contains(c.IgnoreFields, fi.RawName) {
			continue
		}
		// remove unsigned
		fi.Type = strings.TrimPrefix(fi.Type, unsignedPrefix)
		// convert time.Time to int64
		if fi.Type == util.GoTime {
			fi.Type = util.GoInt64
		}
		if fi.Type == util.GoBool && fi.Enums == boolTypeEnums {
			// trim bool comment
			fi.Comment = convertComment(fi.Comment, true)
		} else if fi.Type == util.GoInt {
			if fi.Enums != "" {
				// convert enums int to int32
				fi.Type = util.GoInt32
			} else {
				// convert int to int64 adapted to protobuf
				fi.Type = util.GoInt64
			}
		}
		// get id info
		if fi.IsPrimaryKey {
			gc.IdName = fi.Name
			gc.IdType = fi.Type
			gc.IdRawName = fi.RawName
			if fi.Comment != "" {
				gc.IdComment = fi.Comment
			}
		}
		fields = append(fields, fi)
	}
	gc.StructFields = fields

	return gc
}
