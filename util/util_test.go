package util

import (
	"testing"
)

func TestConvertDataType(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		cc          CMDConfig
		expectation string
	}{
		{ColumnInfo{DataType: "tinyint", IsNullable: true, IsUnsigned: false},
			CMDConfig{EnableSqlNull: false, EnableGureguNull: false}, "int32",
		},
		{ColumnInfo{DataType: "tinyint", IsNullable: false, IsUnsigned: true},
			CMDConfig{EnableSqlNull: false, EnableGureguNull: false}, "uint32",
		},
		{ColumnInfo{DataType: "tinyint", IsNullable: true, IsUnsigned: false},
			CMDConfig{EnableSqlNull: true, EnableGureguNull: false}, "sql.NullInt32",
		},
		{ColumnInfo{DataType: "tinyint", IsNullable: true, IsUnsigned: true},
			CMDConfig{EnableSqlNull: true, EnableGureguNull: true}, "null.Int",
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
		expectation string
	}{
		{"api", "API"},
		{"url", "URL"},
		{"test", "Test"},
		{"user_name", "UserName"},
		{"USER_PASSWORD", "UserPassword"},
		{"测试", "测试"},
		{"用户_名称", "用户名称"},
	}

	for _, c := range cases {
		output := convertName(c.input)
		if output != c.expectation {
			t.Errorf("convertName failed, input:%s, expectation:%s, output:%s",
				c.input, c.expectation, output)
		}
	}
}

func TestGetJsonTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id"}, "json:\"id\""},
		{ColumnInfo{Name: "name"}, "json:\"name\""},
	}

	for _, c := range cases {
		output := getJsonTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getJsonTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetXmlTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id"}, "xml:\"id\""},
		{ColumnInfo{Name: "name"}, "xml:\"name\""},
	}

	for _, c := range cases {
		output := getXmlTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getXmlTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetGormTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id", Type: "bigint(20)", IsPrimaryKey: true,
			IsAutoIncrement: true, IsNullable: false, Default: "", Comment: "用户id"},
			"gorm:\"primary_key;column:id;type:bigint(20) auto_increment;comment:'用户id'\""},
		{ColumnInfo{Name: "name", Type: "varchar(255)", IsPrimaryKey: false,
			IsAutoIncrement: false, IsNullable: false, Default: "user", Comment: "用户名称",
			Indexes: []*IndexInfo{{Name: "name_index"}, {Name: "name_email_index"}}},
			"gorm:\"column:name;type:varchar(255);not null;index:name_index,name_email_index;default:'user';comment:'用户名称'\""},
		{ColumnInfo{Name: "email", Type: "varchar(255)", IsPrimaryKey: false,
			IsAutoIncrement: false, IsNullable: false, Default: "email", Comment: "用户邮箱",
			Indexes:       []*IndexInfo{{Name: "name_email_index"}},
			UniqueIndexes: []*IndexInfo{{Name: "email_index"}}},
			"gorm:\"column:email;type:varchar(255);not null;index:name_email_index;unique_index:email_index;default:'email';comment:'用户邮箱'\""},
	}

	for _, c := range cases {
		output := getGormTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getGormTag failed, expectation:%s, output:%s",
				c.expectation, output)
		}
	}
}

func TestGetXormTag(t *testing.T) {
	cases := []struct {
		ci          ColumnInfo
		expectation string
	}{
		{ColumnInfo{Name: "id", Type: "bigint(20)", IsPrimaryKey: true,
			IsAutoIncrement: true, IsNullable: false, Default: "", Comment: "用户id"},
			"xorm:\"pk autoincr bigint(20) 'id' comment('用户id')\""},
		{ColumnInfo{Name: "name", Type: "varchar(255)", IsPrimaryKey: false,
			IsAutoIncrement: false, IsNullable: false, Default: "user", Comment: "用户名称",
			Indexes: []*IndexInfo{{Name: "name_index"}, {Name: "name_email_index"}}},
			"xorm:\"varchar(255) 'name' notnull index(name_index) index(name_email_index) default('user') comment('用户名称')\""},
		{ColumnInfo{Name: "email", Type: "varchar(255)", IsPrimaryKey: false,
			IsAutoIncrement: false, IsNullable: false, Default: "email", Comment: "用户邮箱",
			Indexes:       []*IndexInfo{{Name: "name_email_index"}},
			UniqueIndexes: []*IndexInfo{{Name: "email_index"}}},
			"xorm:\"varchar(255) 'email' notnull index(name_email_index) unique(email_index) default('email') comment('用户邮箱')\""},
	}

	for _, c := range cases {
		output := getXormTag(&c.ci)
		if output != c.expectation {
			t.Errorf("getXormTag failed, expectation:%s, output:%s",
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
