package api

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"

	"github.com/sliveryou/grom/util"
)

// Config represents the config of the generated api.
type Config struct {
	util.DBConfig

	StructName       string   `json:"struct_name,omitempty"` // camel
	RouteName        string   `json:"route_name,omitempty"`  // snake
	TableComment     string   `json:"-"`
	EnableInitialism bool     `json:"enable_initialism"`
	IgnoreFields     []string `json:"ignore_fields"`

	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Author  string `json:"author"`
	Email   string `json:"email"`
	Version string `json:"version"`

	ServiceName  string `json:"service_name"` // camel
	RoutePrefix  string `json:"route_prefix"` // lower
	GroupPrefix  string `json:"group_prefix"` // lower
	RouteStyle   string `json:"route_style"`  // one of [snake, kebab], default is kebab
	EnablePlural bool   `json:"enable_plural"`
	EnableModel  bool   `json:"enable_model"`
}

// GetCmdConfig gets the *util.CmdConfig.
func (c *Config) GetCmdConfig() *util.CmdConfig {
	return &util.CmdConfig{
		DBConfig:           c.DBConfig,
		PackageName:        "model",
		StructName:         c.StructName,
		EnableInitialism:   c.EnableInitialism,
		EnableFieldComment: true,
		EnableJSONTag:      true,
		EnableGormV2Tag:    true,
		DisableUnsigned:    true,
	}
}

// UpdateBy Config updated by util.CmdConfig.
func (c *Config) UpdateBy(cc *util.CmdConfig) {
	c.StructName = cc.StructName
	c.TableComment = cc.TableComment
}

// GetDelimiter gets the route delimiter.
func (c *Config) GetDelimiter() uint8 {
	var delimiter uint8 = '-'
	if c.RouteStyle == RouteStyleSnake {
		delimiter = '_'
	}

	return delimiter
}

// ProjectConfig represents the config of the generated project.
type ProjectConfig struct {
	Config
	Dir                   string   `json:"dir"`
	TablePrefix           string   `json:"table_prefix"`
	Tables                []string `json:"tables"`
	EnableTrimTablePrefix bool     `json:"enable_trim_table_prefix"`
}

// Check checks whether pc is valid.
func (pc *ProjectConfig) Check() error {
	if pc.Host == "" || pc.Port < 1 ||
		pc.User == "" || pc.Database == "" {
		return errDBConfig
	}
	if pc.ServiceName == "" {
		return errEmptyServiceName
	}
	if pc.Dir == "" {
		return errEmptyDir
	}
	if len(pc.Tables) == 0 {
		return errNoTables
	}
	if pc.RouteStyle != RouteStyleSnake &&
		pc.RouteStyle != RouteStyleKebab {
		pc.RouteStyle = RouteStyleKebab
	}

	return nil
}

// StructField represents the field of the generated model structure.
type StructField struct {
	Name         string
	Type         string
	Comment      string
	RawName      string
	RawType      string
	DataType     string
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
		DataType:     sf.DataType,
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
	IdNamePlural    string
	IdType          string
	IdComment       string
	IdRawName       string
	IdRawNamePlural string
	RouteName       string // snake or kebab
	ModelName       string // camel
	GroupName       string // lower
	StructFields    []StructField
}

func getGenerateConfig(c Config, fs []*util.StructField) generateConfig {
	gc := generateConfig{
		IdComment: c.TableComment + defaultIdComment,
		RouteName: strcase.ToDelimited(c.StructName, c.GetDelimiter()),
		ModelName: c.StructName,
		GroupName: strings.ToLower(c.StructName),
	}
	if c.RouteName != "" {
		gc.RouteName = c.RouteName
		gc.ModelName = strcase.ToCamel(c.RouteName)
		gc.GroupName = strings.ToLower(strcase.ToCamel(c.RouteName))
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
		// convert json to map
		if fi.DataType == dataTypeJSON {
			fi.Type = dataTypeMap
		} else if fi.Type == util.GoTime {
			// convert time.Time to int64
			fi.Type = util.GoInt64
		} else if fi.Type == util.GoBool && fi.Enums == boolTypeEnums {
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
			gc.IdNamePlural = inflection.Plural(gc.IdName)
			gc.IdType = fi.Type
			gc.IdRawName = fi.RawName
			gc.IdRawNamePlural = inflection.Plural(gc.IdRawName)
			if fi.Comment != "" {
				gc.IdComment = fi.Comment
			}
		}
		fields = append(fields, fi)
	}
	gc.StructFields = fields

	return gc
}

func cloneStructFields(cc *util.CmdConfig, fs []*util.StructField) []*util.StructField {
	cloneFs := make([]*util.StructField, 0, len(fs))

	for _, f := range fs {
		cloneF := *f
		enums := getEnums(f.Comment)
		cloneF.Name = initialismsReplacer.Replace(cloneF.Name)
		cloneF.Type = strings.TrimPrefix(cloneF.Type, unsignedPrefix)

		if cloneF.Type == util.GoTime {
			if cloneF.RawName == deleteAt {
				cloneF.Type = gormDeleteAt
				cc.EnableGROM = true
			} else {
				cloneF.Type = util.GoPointerTime
			}
		} else if cloneF.Type == util.GoBool && enums == boolTypeEnums {
			// trim bool comment
			cloneF.Comment = convertComment(cloneF.Comment, true)
		} else if cloneF.Type == util.GoInt {
			if enums != "" {
				// convert enums int to int32
				cloneF.Type = util.GoInt32
			} else {
				// convert int to int64 adapted to protobuf
				cloneF.Type = util.GoInt64
			}
		}
		if cloneF.DataType == dataTypeJSON {
			cloneF.Type = dataTypesJSON
			cc.EnableDataTypes = true
		}
		if cloneF.Default != "" {
			cloneF.Type = toPointer(cloneF.Type)
		}

		cloneFs = append(cloneFs, &cloneF)
	}

	return cloneFs
}
