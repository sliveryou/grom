package util

import (
	"testing"
)

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
			t.Errorf("convertName failed, input:%s, expectation:%s, output:%s", c.input, c.expectation, output)
		}
	}
}
