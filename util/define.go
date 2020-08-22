package util

const (
	// MySQLDriverName represents the mysql driver name.
	MySQLDriverName = "mysql"
)

// CMDConfig represents the config of the running grom command line.
type CMDConfig struct {
	DBConfig
	PackageName        string `json:"package_name"`
	StructName         string `json:"struct_name"`
	EnableFieldComment bool   `json:"enable_field_comment"`
	EnableSqlNull      bool   `json:"enable_sql_null"`
	EnableGureguNull   bool   `json:"enable_guregu_null"`
	EnableJsonTag      bool   `json:"enable_json_tag"`
	EnableXmlTag       bool   `json:"enable_xml_tag"`
	EnableGormTag      bool   `json:"enable_gorm_tag"`
	EnableXormTag      bool   `json:"enable_xorm_tag"`
	EnableBeegoTag     bool   `json:"enable_beego_tag"`
	EnableGoroseTag    bool   `json:"enable_gorose_tag"`
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
