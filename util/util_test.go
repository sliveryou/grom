package util

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConvertTable(t *testing.T) {
	config := CmdConfig{
		DBConfig: DBConfig{
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
		EnableSQLNull:      false,
		EnableGureguNull:   false,
		EnableJSONTag:      true,
		EnableXMLTag:       false,
		EnableGormTag:      true,
		EnableXormTag:      false,
		EnableBeegoTag:     false,
		EnableGoroseTag:    false,
	}

	b, err := json.MarshalIndent(&config, "", "    ")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("\n" + string(b))
	}

	// s, err := ConvertTable(config)
	// if err != nil {
	// 	t.Error(err)
	// } else {
	// 	t.Log("\n" + s)
	// }
}

func TestConvertDataType(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		cc          CmdConfig
		expectation string
	}{
		{
			ColumnInfo{DataType: "tinyint", IsNullable: true, IsUnsigned: false},
			CmdConfig{EnableSQLNull: false, EnableGureguNull: false},
			"int32",
		},
		{
			ColumnInfo{DataType: "tinyint", IsNullable: false, IsUnsigned: true},
			CmdConfig{EnableSQLNull: false, EnableGureguNull: false},
			"uint32",
		},
		{
			ColumnInfo{DataType: "tinyint", IsNullable: true, IsUnsigned: false},
			CmdConfig{EnableSQLNull: true, EnableGureguNull: false},
			"sql.NullInt32",
		},
		{
			ColumnInfo{DataType: "tinyint", IsNullable: true, IsUnsigned: true},
			CmdConfig{EnableSQLNull: true, EnableGureguNull: true},
			"null.Int",
		},
	}

	for _, c := range cases {
		output := convertDataType(&c.ci, &c.cc)
		if output != c.expectation {
			t.Errorf("convertDataType failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestConvertName(t *testing.T) {
	cases := []struct {
		input       string
		enable      bool
		expectation string
	}{
		{"api", true, "API"},
		{"url", true, "URL"},
		{"api", false, "Api"},
		{"url", false, "Url"},
		{"test", true, "Test"},
		{"user_name", true, "UserName"},
		{"USER_PASSWORD", true, "UserPassword"},
		{"测试", true, "测试"},
		{"用户_名称", true, "用户名称"},
		{"user_name", false, "UserName"},
		{"USER_PASSWORD", false, "UserPassword"},
		{"测试", false, "测试"},
		{"用户_名称", false, "用户名称"},
	}

	for _, c := range cases {
		output := convertName(c.input, c.enable)
		if output != c.expectation {
			t.Errorf("convertName failed, input:%s, enable:%v, expectation:%s, output:%s",
				c.input, c.enable, c.expectation, output)
		}
	}
}

func TestGetJSONTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id"}, "json:\"id\""},
		{ColumnInfo{Name: "name"}, "json:\"name\""},
	}

	for _, c := range cases {
		output := getJSONTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getJSONTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetXMLTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id"}, "xml:\"id\""},
		{ColumnInfo{Name: "name"}, "xml:\"name\""},
	}

	for _, c := range cases {
		output := getXMLTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getXMLTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetGormTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{
			ColumnInfo{
				Name: "id", Type: "bigint(20)", IsPrimaryKey: true,
				IsAutoIncrement: true, IsNullable: false, Default: "", Comment: "用户id",
			},
			"gorm:\"primary_key;column:id;type:bigint(20) auto_increment;comment:'用户id'\"",
		},
		{
			ColumnInfo{
				Name: "name", Type: "varchar(255)", IsPrimaryKey: false,
				IsAutoIncrement: false, IsNullable: false, Default: "user", Comment: "用户名称",
				Indexes: []*IndexInfo{{Name: "name_index"}, {Name: "name_email_index"}},
			},
			"gorm:\"column:name;type:varchar(255);not null;index:name_index,name_email_index;default:'user';comment:'用户名称'\"",
		},
		{
			ColumnInfo{
				Name: "email", Type: "varchar(255)", IsPrimaryKey: false,
				IsAutoIncrement: false, IsNullable: false, Default: "email", Comment: "用户邮箱",
				Indexes:       []*IndexInfo{{Name: "name_email_index"}},
				UniqueIndexes: []*IndexInfo{{Name: "email_index"}},
			},
			"gorm:\"column:email;type:varchar(255);not null;index:name_email_index;unique_index:email_index;default:'email';comment:'用户邮箱'\"",
		},
	}

	for _, c := range cases {
		output := getGormTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getGormTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetGormV2Tag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{
			ColumnInfo{
				Name: "id", Type: "bigint(20)", IsPrimaryKey: true,
				IsAutoIncrement: true, IsNullable: false, Default: "", Comment: "用户id",
			},
			"gorm:\"primaryKey;autoIncrement;column:id;comment:用户id\"",
		},
		{
			ColumnInfo{
				Name: "name", Type: "varchar(255)", IsPrimaryKey: false,
				IsAutoIncrement: false, IsNullable: false, Default: "user", Comment: "用户名称",
				Indexes: []*IndexInfo{{Name: "name_index"}, {Name: "name_email_index"}},
			},
			"gorm:\"column:name;type:varchar(255);not null;index:name_index,name_email_index;default:user;comment:用户名称\"",
		},
		{
			ColumnInfo{
				Name: "email", Type: "varchar(255)", IsPrimaryKey: false,
				IsAutoIncrement: false, IsNullable: false, Default: "email", Comment: "用户邮箱",
				Indexes:       []*IndexInfo{{Name: "name_email_index"}},
				UniqueIndexes: []*IndexInfo{{Name: "email_index"}},
			},
			"gorm:\"column:email;type:varchar(255);not null;index:name_email_index;uniqueIndex:email_index;default:email;comment:用户邮箱\"",
		},
	}

	for _, c := range cases {
		output := getGormV2Tag(&c.ci)
		if output != c.expectation {
			t.Errorf("getGormV2Tag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetXormTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{
			ColumnInfo{
				Name: "id", Type: "bigint(20)", IsPrimaryKey: true,
				IsAutoIncrement: true, IsNullable: false, Default: "", Comment: "用户id",
			},
			"xorm:\"pk autoincr bigint(20) 'id' comment('用户id')\"",
		},
		{
			ColumnInfo{
				Name: "name", Type: "varchar(255)", IsPrimaryKey: false,
				IsAutoIncrement: false, IsNullable: false, Default: "user", Comment: "用户名称",
				Indexes: []*IndexInfo{{Name: "name_index"}, {Name: "name_email_index"}},
			},
			"xorm:\"varchar(255) 'name' notnull index(name_index) index(name_email_index) default('user') comment('用户名称')\"",
		},
		{
			ColumnInfo{
				Name: "email", Type: "varchar(255)", IsPrimaryKey: false,
				IsAutoIncrement: false, IsNullable: false, Default: "email", Comment: "用户邮箱",
				Indexes:       []*IndexInfo{{Name: "name_email_index"}},
				UniqueIndexes: []*IndexInfo{{Name: "email_index"}},
			},
			"xorm:\"varchar(255) 'email' notnull index(name_email_index) unique(email_index) default('email') comment('用户邮箱')\"",
		},
	}

	for _, c := range cases {
		output := getXormTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getXormTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetBeegoTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{
			ColumnInfo{
				Name: "id", Type: "bigint(20)", DataType: "bigint", Precision: 19,
				Scale: 0, Length: 0, IsPrimaryKey: true, IsAutoIncrement: true,
				IsNullable: false, Default: "", Comment: "用户id",
			},
			"orm:\"pk;auto;column(id);type(bigint);size(20);description(用户id)\"",
		},
		{
			ColumnInfo{
				Name: "name", Type: "varchar(255)", DataType: "varchar", Precision: 0,
				Scale: 0, Length: 255, IsPrimaryKey: false, IsAutoIncrement: false,
				IsNullable: false, Default: "user", Comment: "用户名称",
			},
			"orm:\"column(name);type(varchar);size(255);default(user);description(用户名称)\"",
		},
		{
			ColumnInfo{
				Name: "email", Type: "varchar(255)", DataType: "varchar", Precision: 0,
				Scale: 0, Length: 255, IsPrimaryKey: false, IsAutoIncrement: false,
				IsNullable: false, Default: "email", Comment: "用户邮箱",
			},
			"orm:\"column(email);type(varchar);size(255);default(email);description(用户邮箱)\"",
		},
	}

	for _, c := range cases {
		output := getBeegoTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getBeegoTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetGoroseTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id"}, "gorose:\"id\""},
		{ColumnInfo{Name: "name"}, "gorose:\"name\""},
	}

	for _, c := range cases {
		output := getGoroseTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getGoroseTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetBeegoType(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{
			ColumnInfo{DataType: "decimal", Length: 0, Precision: 10, Scale: 2, IsUnsigned: true},
			";type(decimal unsigned);digits(10);decimals(2)",
		},
		{
			ColumnInfo{DataType: "bigint", Length: 0, Precision: 19, Scale: 0, IsUnsigned: false},
			";type(bigint);size(20)",
		},
		{
			ColumnInfo{DataType: "varchar", Length: 255, Precision: 0, Scale: 0, IsUnsigned: false},
			";type(varchar);size(255)",
		},
		{
			ColumnInfo{DataType: "longtext", Length: 0, Precision: 0, Scale: 0, IsUnsigned: false},
			";type(text)",
		},
		{
			ColumnInfo{DataType: "date", Length: 0, Precision: 0, Scale: 0, IsUnsigned: false},
			";type(date)",
		},
		{
			ColumnInfo{DataType: "time", Length: 0, Precision: 0, Scale: 0, IsUnsigned: false},
			";type(datetime)",
		},
		{
			ColumnInfo{DataType: "longblob", Length: 0, Precision: 0, Scale: 0, IsUnsigned: false},
			";type(longblob)",
		},
		{
			ColumnInfo{DataType: "longblob", Length: 0, Precision: 0, Scale: 0, IsUnsigned: false},
			";type(longblob)",
		},
	}

	for _, c := range cases {
		output := getBeegoType(&c.ci)
		if output != c.expectation {
			t.Errorf("getBeegoType failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestUniqueStrings(t *testing.T) {
	cases := []struct {
		input       []string
		expectation []string
	}{
		{[]string{"a", "a", "b", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "a", "c", "b"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, c := range cases {
		output := uniqueStrings(c.input)
		if !reflect.DeepEqual(output, c.expectation) {
			t.Errorf("uniqueStrings failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}
