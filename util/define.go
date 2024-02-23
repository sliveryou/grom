package util

import (
	"database/sql"
	"sync"
)

// Global variables.
var (
	db      *sql.DB
	dbMutex sync.Mutex

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
	GureguNullString = "null.String"
	GureguNullInt    = "null.Int"
	GureguNullFloat  = "null.Float"
	GureguNullBool   = "null.Bool"
	GureguNullTime   = "null.Time"

	SQLNullString  = "sql.NullString"
	SQLNullInt32   = "sql.NullInt32"
	SQLNullInt64   = "sql.NullInt64"
	SQLNullFloat64 = "sql.NullFloat64"
	SQLNullBool    = "sql.NullBool"
	SQLNullTime    = "sql.NullTime"

	GoString      = "string"
	GoBytes       = "[]byte"
	GoInt         = "int"
	GoUint        = "uint"
	GoInt32       = "int32"
	GoUint32      = "uint32"
	GoInt64       = "int64"
	GoUint64      = "uint64"
	GoFloat32     = "float32"
	GoFloat64     = "float64"
	GoBool        = "bool"
	GoTime        = "time.Time"
	GoPointerTime = "*time.Time"
)

// CmdConfig represents the config of the running grom command line.
type CmdConfig struct {
	DBConfig
	PackageName        string   `json:"package_name"`
	StructName         string   `json:"struct_name"`
	EnableInitialism   bool     `json:"enable_initialism"`
	EnableFieldComment bool     `json:"enable_field_comment"`
	EnableSQLNull      bool     `json:"enable_sql_null"`
	EnableGureguNull   bool     `json:"enable_guregu_null"`
	EnableJSONTag      bool     `json:"enable_json_tag"`
	EnableXMLTag       bool     `json:"enable_xml_tag"`
	EnableGormTag      bool     `json:"enable_gorm_tag"`
	EnableXormTag      bool     `json:"enable_xorm_tag"`
	EnableBeegoTag     bool     `json:"enable_beego_tag"`
	EnableGoroseTag    bool     `json:"enable_gorose_tag"`
	EnableGormV2Tag    bool     `json:"enable_gorm_v2_tag"`
	DisableUnsigned    bool     `json:"disable_unsigned"`
	EnableGoTime       bool     `json:"-"`
	EnableGROM         bool     `json:"-"`
	EnableDataTypes    bool     `json:"-"`
	TableComment       string   `json:"-"`
	TableIndexes       []string `json:"-"`
	TableUniques       []string `json:"-"`
}

// DBConfig represents the config of the connected database.
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Table    string `json:"table,omitempty"`
}

// StructField represents the field of the generated model structure.
type StructField struct {
	Name         string
	Type         string
	Tag          string
	Comment      string
	RawName      string
	DataType     string
	Default      string
	IsPrimaryKey bool
	IsNullable   bool
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
	Another         string       `mysql:"-"`
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
