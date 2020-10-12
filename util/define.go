package util

import (
	"database/sql"
)

// Global variables.
var (
	db *sql.DB

	columnInfos []*ColumnInfo
	indexInfos  []*IndexInfo

	tableIndexes []string
	tableUniques []string

	commonInitialisms = map[string]struct{}{
		"ACL":   {},
		"API":   {},
		"ASCII": {},
		"CPU":   {},
		"CSS":   {},
		"DNS":   {},
		"EOF":   {},
		"GUID":  {},
		"HTML":  {},
		"HTTP":  {},
		"HTTPS": {},
		"ID":    {},
		"IP":    {},
		"JSON":  {},
		"LHS":   {},
		"QPS":   {},
		"RAM":   {},
		"RHS":   {},
		"RPC":   {},
		"SLA":   {},
		"SMTP":  {},
		"SQL":   {},
		"SSH":   {},
		"TCP":   {},
		"TLS":   {},
		"TTL":   {},
		"UDP":   {},
		"UI":    {},
		"UID":   {},
		"UUID":  {},
		"URI":   {},
		"URL":   {},
		"UTF8":  {},
		"VM":    {},
		"XML":   {},
		"XMPP":  {},
		"XSRF":  {},
		"XSS":   {},
	}
)

const (
	// MySQLDriverName represents the mysql driver name.
	MySQLDriverName = "mysql"
)

const (
	indexUnique = 0
	indexNormal = 1
)

// Global data type constants.
const (
	gureguNullString = "null.String"
	gureguNullInt    = "null.Int"
	gureguNullFloat  = "null.Float"
	gureguNullBool   = "null.Bool"
	gureguNullTime   = "null.Time"

	sqlNullString  = "sql.NullString"
	sqlNullInt32   = "sql.NullInt32"
	sqlNullInt64   = "sql.NullInt64"
	sqlNullFloat64 = "sql.NullFloat64"
	sqlNullBool    = "sql.NullBool"
	sqlNullTime    = "sql.NullTime"

	goString  = "string"
	goBytes   = "[]byte"
	goInt     = "int"
	goUint    = "uint"
	goInt32   = "int32"
	goUint32  = "uint32"
	goInt64   = "int64"
	goUint64  = "uint64"
	goFloat32 = "float32"
	goFloat64 = "float64"
	goBool    = "bool"
	goTime    = "time.Time"
)

// CMDConfig represents the config of the running grom command line.
type CMDConfig struct {
	DBConfig
	PackageName        string `json:"package_name"`
	StructName         string `json:"struct_name"`
	EnableInitialism   bool   `json:"enable_initialism"`
	EnableFieldComment bool   `json:"enable_field_comment"`
	EnableSqlNull      bool   `json:"enable_sql_null"`
	EnableGureguNull   bool   `json:"enable_guregu_null"`
	EnableJsonTag      bool   `json:"enable_json_tag"`
	EnableXmlTag       bool   `json:"enable_xml_tag"`
	EnableGormTag      bool   `json:"enable_gorm_tag"`
	EnableXormTag      bool   `json:"enable_xorm_tag"`
	EnableBeegoTag     bool   `json:"enable_beego_tag"`
	EnableGoroseTag    bool   `json:"enable_gorose_tag"`
	EnableGormV2Tag    bool   `json:"enable_gorm_v2_tag"`
}

// DBConfig represents the config of the connected database.
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

// StructField represents the field of the generated model structure.
type StructField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

// ColumnInfo represents the information of the column.
type ColumnInfo struct {
	Name            string       `mysql:"COLUMN_NAME"`
	DataType        string       `mysql:"DATA_TYPE"`
	Type            string       `mysql:"COLUMN_TYPE"`
	Default         string       `mysql:"COLUMN_DEFAULT"`
	Comment         string       `mysql:"COLUMN_COMMENT"`
	Length          int64        `mysql:"CHARACTER_MAXIMUM_LENGTH"`
	Precision       int64        `mysql:"NUMERIC_PRECISION"`
	Scale           int64        `mysql:"NUMERIC_SCALE"`
	Position        int          `mysql:"ORDINAL_POSITION"`
	IsPrimaryKey    bool         `mysql:"COLUMN_KEY"`
	IsAutoIncrement bool         `mysql:"EXTRA"`
	IsUnsigned      bool         `mysql:"COLUMN_TYPE"`
	IsNullable      bool         `mysql:"IS_NULLABLE"`
	Indexes         []*IndexInfo `mysql:"-"`
	UniqueIndexes   []*IndexInfo `mysql:"-"`
}

// IndexInfo represents the information of the index.
type IndexInfo struct {
	Name       string `mysql:"INDEX_NAME"`
	ColumnName string `mysql:"COLUMN_NAME"`
	Comment    string `mysql:"INDEX_COMMENT"`
	Sequence   int    `mysql:"SEQ_IN_INDEX"`
	IsUnique   bool   `mysql:"NON_UNIQUE"`
}
