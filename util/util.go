package util

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// ConvertTable converts mysql table fields to golang model structure by command config.
func ConvertTable(cc CMDConfig) (string, error) {
	defer CloseDB()

	fields, err := GetFields(&cc)
	if err != nil {
		return "", err
	}

	return generateCode(&cc, fields)
}

// GetFields gets golang structure fields converted by mysql table fields.
func GetFields(cc *CMDConfig) ([]*StructField, error) {
	if cc.PackageName == "" {
		cc.PackageName = "model"
	}
	if cc.StructName == "" {
		cc.StructName = convertName(cc.Table, cc.EnableInitialism)
	}

	comment, err := getTableComment(cc)
	if err != nil {
		return nil, errors.WithMessage(err, "getTableComment err")
	}
	cc.TableComment = comment

	cis, err := getColumnInfos(cc)
	if err != nil {
		return nil, errors.WithMessage(err, "getColumnInfos err")
	}

	var fields []*StructField
	for i := range cis {
		ci := cis[i]
		var tags []string

		if cc.EnableJsonTag {
			tags = append(tags, getJsonTag(ci))
		}
		if cc.EnableXmlTag {
			tags = append(tags, getXmlTag(ci))
		}
		if cc.EnableGormTag {
			tags = append(tags, getGormTag(ci))
		}
		if cc.EnableXormTag {
			tags = append(tags, getXormTag(ci))
		}
		if cc.EnableBeegoTag {
			tags = append(tags, getBeegoTag(ci))
		}
		if cc.EnableGoroseTag {
			tags = append(tags, getGoroseTag(ci))
		}
		if cc.EnableGormV2Tag && !cc.EnableGormTag {
			tags = append(tags, getGormV2Tag(ci))
		}

		field := StructField{
			Name:         convertName(ci.Name, cc.EnableInitialism),
			Type:         convertDataType(ci, cc),
			Comment:      ci.Comment,
			RawName:      ci.Name,
			Default:      ci.Default,
			IsPrimaryKey: ci.IsPrimaryKey,
			IsNullable:   ci.IsNullable,
		}
		if len(tags) > 0 {
			field.Tag = fmt.Sprintf("`%s`", strings.Join(removeEmpty(tags), " "))
		}
		if field.Type == GoTime {
			cc.EnableGoTime = true
		}
		fields = append(fields, &field)
	}

	return fields, nil
}

// convertDataType converts the mysql data type to golang data type.
func convertDataType(ci *ColumnInfo, cc *CMDConfig) string {
	switch ci.DataType {
	case "tinyint", "smallint", "mediumint":
		isBool := false
		if strings.Contains(ci.Type, "tinyint(1)") {
			isBool = true
		}
		if ci.IsNullable {
			if cc.EnableGureguNull {
				if isBool {
					return GureguNullBool
				}
				return GureguNullInt
			} else if cc.EnableSqlNull {
				if isBool {
					return SqlNullBool
				}
				return SqlNullInt32
			}
		}
		if ci.IsUnsigned {
			if isBool {
				return GoBool
			}
			return GoUint32
		}
		if isBool {
			return GoBool
		}
		return GoInt32
	case "int", "integer":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return GureguNullInt
			} else if cc.EnableSqlNull {
				return SqlNullInt64
			}
		}
		if ci.IsUnsigned {
			return GoUint
		}
		return GoInt
	case "bigint":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return GureguNullInt
			} else if cc.EnableSqlNull {
				return SqlNullInt64
			}
		}
		if ci.IsUnsigned {
			return GoUint64
		}
		return GoInt64
	case "json", "enum", "set", "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return GureguNullString
			} else if cc.EnableSqlNull {
				return SqlNullString
			}
		}
		return GoString
	case "year", "date", "datetime", "time", "timestamp":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return GureguNullTime
			} else if cc.EnableSqlNull {
				return SqlNullTime
			}
		}
		return GoTime
	case "float":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return GureguNullFloat
			} else if cc.EnableSqlNull {
				return SqlNullFloat64
			}
		}
		return GoFloat32
	case "double", "real", "decimal", "numeric":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return GureguNullFloat
			} else if cc.EnableSqlNull {
				return SqlNullFloat64
			}
		}
		return GoFloat64
	case "bit", "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return GoBytes
	default:
		return "unknown"
	}
}

// convertName converts the name to camel case name.
func convertName(name string, enableInitialism ...bool) string {
	if name == "" {
		return ""
	}

	enable := false
	if len(enableInitialism) != 0 {
		enable = enableInitialism[0]
	}

	var cn string
	s := strings.Split(name, "_")

	for _, v := range s {
		upperV := strings.ToUpper(v)
		if _, ok := commonInitialisms[upperV]; ok && enable {
			cn += upperV
		} else {
			if runesV := []rune(v); len(runesV) > 0 {
				for i, r := range runesV {
					if i == 0 {
						runesV[i] = unicode.ToUpper(r)
					} else {
						runesV[i] = unicode.ToLower(r)
					}
				}
				cn += string(runesV)
			}
		}
	}

	return cn
}

// getJsonTag returns the tag string of json.
func getJsonTag(ci *ColumnInfo) string {
	return fmt.Sprintf("json:%q", ci.Name)
}

// getXmlTag returns the tag string of xml.
func getXmlTag(ci *ColumnInfo) string {
	return fmt.Sprintf("xml:%q", ci.Name)
}

// getGormTag returns the tag string of gorm.
func getGormTag(ci *ColumnInfo) string {
	return generateTag(ci, gormTplName)
}

// getXormTag returns the tag string of xorm.
func getXormTag(ci *ColumnInfo) string {
	return generateTag(ci, xormTplName)
}

// getBeegoTag returns the tag string of beego orm.
func getBeegoTag(ci *ColumnInfo) string {
	return generateTag(ci, beegoTplName)
}

// getGoroseTag returns the tag string of gorose.
func getGoroseTag(ci *ColumnInfo) string {
	return fmt.Sprintf("gorose:%q", ci.Name)
}

// getGormV2Tag returns the tag string of gorm v2.
func getGormV2Tag(ci *ColumnInfo) string {
	return generateTag(ci, gormV2TplName)
}

// getBeegoType returns the type tag string of beego orm.
func getBeegoType(ci *ColumnInfo) string {
	sign := ""
	if ci.IsUnsigned {
		sign = " unsigned"
	}

	switch ci.DataType {
	case "float", "double", "real", "decimal", "numeric":
		return fmt.Sprintf(";type(%s%s);digits(%d);decimals(%d)", ci.DataType, sign, ci.Precision, ci.Scale)
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		return fmt.Sprintf(";type(%s%s);size(%d)", ci.DataType, sign, ci.Precision+1)
	case "date", "datetime":
		return fmt.Sprintf(";type(%s)", ci.DataType)
	case "year", "time", "timestamp":
		return ";type(datetime)"
	case "bit", "binary", "varbinary", "char", "varchar":
		return fmt.Sprintf(";type(%s);size(%d)", ci.DataType, ci.Length)
	case "tinytext", "text", "mediumtext", "longtext":
		return ";type(text)"
	case "json", "enum", "set", "tinyblob", "blob", "mediumblob", "longblob":
		return fmt.Sprintf(";type(%s)", ci.DataType)
	default:
		return fmt.Sprintf(";type(%s)", ci.DataType)
	}
}

// removeEmpty remove empty fields.
func removeEmpty(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, field := range slice {
		if field != "" {
			result = append(result, field)
		}
	}

	return result
}
