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
