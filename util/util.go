package util

import (
	"strings"
	"unicode"
)

// convertDataType converts the mysql data type to golang data type.
func convertDataType(ci *ColumnInfo, cc *CMDConfig) string {
	switch ci.DataType {
	case "tinyint", "smallint", "mediumint":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullInt
			} else if cc.EnableSqlNull {
				return sqlNullInt32
			}
		}
		if ci.IsUnsigned {
			return goUint32
		}
		return goInt32
	case "int", "integer":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullInt
			} else if cc.EnableSqlNull {
				return sqlNullInt64
			}
		}
		if ci.IsUnsigned {
			return goUint
		}
		return goInt
	case "bigint":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullInt
			} else if cc.EnableSqlNull {
				return sqlNullInt64
			}
		}
		if ci.IsUnsigned {
			return goUint64
		}
		return goInt64
	case "json", "enum", "set", "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullString
			} else if cc.EnableSqlNull {
				return sqlNullString
			}
		}
		return goString
	case "year", "date", "datetime", "time", "timestamp":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullTime
			} else if cc.EnableSqlNull {
				return sqlNullTime
			}
		}
		return goTime
	case "float":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullFloat
			} else if cc.EnableSqlNull {
				return sqlNullFloat64
			}
		}
		return goFloat32
	case "double", "real", "decimal", "numeric":
		if ci.IsNullable {
			if cc.EnableGureguNull {
				return gureguNullFloat
			} else if cc.EnableSqlNull {
				return sqlNullFloat64
			}
		}
		return goFloat64
	case "bit", "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return goBytes
	default:
		return "unknown"
	}
}

// convertName converts the name to camel case name.
func convertName(name string) string {
	if name == "" {
		return ""
	}

	var cn string
	s := strings.Split(name, "_")

	for _, v := range s {
		upperV := strings.ToUpper(v)
		if _, ok := abbreviation[upperV]; ok {
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
