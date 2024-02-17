package api

import (
	"testing"

	"github.com/jinzhu/inflection"
	"github.com/stretchr/testify/assert"
)

func Test_initialisms_Replace(t *testing.T) {
	cases := []struct {
		src    string
		expect string
	}{
		{src: "Api", expect: "API"},
		{src: "Id", expect: "ID"},
		{src: "Ip", expect: "IP"},
		{src: "Json", expect: "JSON"},
		{src: "Abc", expect: "Abc"},
		{src: "", expect: ""},
	}

	for _, c := range cases {
		get := initialismsReplacer.Replace(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func Test_getTypeEmptyString(t *testing.T) {
	cases := []struct {
		src    string
		expect string
	}{
		{src: "int", expect: "0"},
		{src: "int64", expect: "0"},
		{src: "float64", expect: "0"},
		{src: "bool", expect: "false"},
		{src: "[]byte", expect: "nil"},
		{src: "*time.Time", expect: "nil"},
		{src: "time.Time", expect: "time.Now()"},
		{src: "string", expect: `""`},
		{src: "map[string]interface{}", expect: "nil"},
		{src: "**string", expect: "nil"},
		{src: "", expect: ""},
	}

	for _, c := range cases {
		get := getTypeEmptyString(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func Test_toPointer(t *testing.T) {
	cases := []struct {
		src    string
		expect string
	}{
		{src: "int", expect: "*int"},
		{src: "*int", expect: "*int"},
		{src: "[]string", expect: "[]string"},
		{src: "map[string]interface{}", expect: "map[string]interface{}"},
		{src: "**string", expect: "**string"},
		{src: "", expect: ""},
	}

	for _, c := range cases {
		get := toPointer(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func Test_isReferenceType(t *testing.T) {
	cases := []struct {
		src    string
		expect bool
	}{
		{src: "[]int", expect: true},
		{src: "*[]int", expect: true},
		{src: "[]string", expect: true},
		{src: "*[]string", expect: true},
		{src: "***[]string", expect: true},
		{src: "map[string]interface{}", expect: true},
		{src: "*map[string]interface{}", expect: true},
		{src: "map[string]int", expect: true},
		{src: "*map[string]int", expect: true},
		{src: "int", expect: false},
		{src: "string", expect: false},
		{src: "", expect: false},
	}

	for _, c := range cases {
		get := isReferenceType(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func Test_toAbbr(t *testing.T) {
	cases := []struct {
		src    string
		expect string
	}{
		{src: "MyName", expect: "mn"},
		{src: "SliverYou", expect: "sy"},
		{src: "HTTPS", expect: "h"},
		{src: "ip", expect: "i"},
		{src: "Json", expect: "j"},
		{src: "Abc", expect: "a"},
		{src: "APIInterfaceDoc", expect: "ad"},
		{src: "API_Interface_Doc", expect: "aid"},
		{src: "my_name_is", expect: "mni"},
		{src: "", expect: ""},
	}

	for _, c := range cases {
		get := toAbbr(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func Test_inflection_Singular(t *testing.T) {
	cases := []struct {
		src    string
		expect string
	}{
		{src: "MyName", expect: "MyName"},
		{src: "APPLES", expect: "APPLE"},
		{src: "companies", expect: "company"},
		{src: "string_numbers", expect: "string_number"},
		{src: "orange", expect: "orange"},
		{src: "id", expect: "id"},
		{src: "", expect: ""},
	}

	for _, c := range cases {
		get := inflection.Singular(c.src)
		assert.Equal(t, c.expect, get)
	}
}
