package util

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

var generator *template.Template

const outTemplate = `
package {{.PackageName}}

{{ if or .EnableSqlNull .EnableGureguNull }}
import (
	{{ if .EnableGureguNull }}
		"gopkg.in/guregu/null.v4"
	{{ else if .EnableSqlNull }}
		"database/sql"
	{{ end }}
)
{{ end }}

type {{ .StructName }} struct {
	{{ range .StructFields -}} 
		{{ .Name }} {{ .Type }} {{ .Tag }} 
		{{- if and $.EnableFieldComment .Comment }}// {{ .Comment }}{{ end }} 
	{{ end -}}
}

{{ if .EnableTableName }}
// TableName returns the table name of the {{ .StructName }} model
func ({{ .ShortStructName }} *{{ .StructName }}) TableName() string {
	return "{{ .Table }}"
}
{{ end }}

{{ if .EnableTableIndex }}
// TableIndex returns the table indexes of the {{ .StructName }} model
func ({{ .ShortStructName }} *{{ .StructName }}) TableIndex() [][]string {
	return [][]string{
		{{ range .TableIndexes -}}
			{ {{ . }} }, 
		{{ end }}
	}
}
{{ end }}

{{ if .EnableTableUnique }}
// TableUnique returns the table unique indexes of the {{ .StructName }} model
func ({{ .ShortStructName }} *{{ .StructName }}) TableUnique() [][]string {
	return [][]string{
		{{ range .TableUniques -}}
			{ {{ . }} }, 
		{{ end }}
	}
}
{{ end }}
`

const gormTemplate = `gorm:"
	{{- if .IsPrimaryKey }}primary_key;{{ end -}}
	column:{{ .Name }};type:{{ .Type }}{{ if .IsAutoIncrement }} auto_increment{{ end }}
	{{- if or .IsNullable .IsPrimaryKey | not }};not null{{ end -}}
	{{- range $i, $v := .Indexes }}
		{{- if eq $i 0 }};index:{{ $v.Name }}{{ else }},{{ $v.Name }}{{ end }}{{ end -}}
	{{- range $i, $v := .UniqueIndexes }}
		{{- if eq $i 0 }};unique_index:{{ $v.Name }}{{ else }},{{ $v.Name }}{{ end }}{{ end -}}
	{{- if .Default }};default:'{{ .Default }}'{{ end -}}
	{{- if .Comment }};comment:'{{ .Comment }}'{{ end -}}
"`

const xormTemplate = `xorm:"
	{{- if .IsPrimaryKey }}pk {{ end -}}
	{{- if .IsAutoIncrement }}autoincr {{ end -}}
	{{ .Type }} '{{ .Name }}'
	{{- if or .IsNullable .IsPrimaryKey | not }} notnull{{ end -}}
	{{- range $i, $v := .Indexes }} index({{ $v.Name }}){{ end -}}
	{{- range $i, $v := .UniqueIndexes }} unique({{ $v.Name }}){{ end -}}
	{{- if .Default }} default('{{ .Default }}'){{ end -}}
	{{- if .Comment }} comment('{{ .Comment }}'){{ end -}}
"`

const beegoTemplate = `orm:"
	{{- if .IsPrimaryKey }}pk;{{ end -}}
	{{- if .IsAutoIncrement }}auto;{{ end -}}
	column({{ .Name }}){{ getBeegoType . }}
	{{- if .IsNullable }};null{{ end -}}
	{{- if .Default }};default({{ .Default }}){{ end -}}
	{{- if .Comment }};description({{ .Comment }}){{ end -}}
"`

func init() {
	var err error
	generator, err = template.New("out").Parse(outTemplate)
	if err != nil {
		fmt.Println("parse out template err:", err)
	}
	generator, err = generator.New("gorm").Parse(gormTemplate)
	if err != nil {
		fmt.Println("parse gorm template err:", err)
	}
	generator, err = generator.New("xorm").Parse(xormTemplate)
	if err != nil {
		fmt.Println("parse xorm template err:", err)
	}
	generator, err = generator.New("beego").Funcs(
		template.FuncMap{"getBeegoType": getBeegoType}).Parse(beegoTemplate)
	if err != nil {
		fmt.Println("parse beego orm template err:", err)
	}
}

// generateCode generates the output code by command config and structure fields.
func generateCode(cc *CMDConfig, fields []*StructField) (string, error) {
	if cc.PackageName == "" {
		cc.PackageName = "model"
	}

	if cc.StructName == "" {
		cc.StructName = convertName(cc.Table)
	}

	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, "out", struct {
		Table              string
		PackageName        string
		StructName         string
		ShortStructName    string
		StructFields       []*StructField
		TableIndexes       []string
		TableUniques       []string
		EnableFieldComment bool
		EnableSqlNull      bool
		EnableGureguNull   bool
		EnableTableName    bool
		EnableTableIndex   bool
		EnableTableUnique  bool
	}{
		Table:              cc.Table,
		PackageName:        cc.PackageName,
		StructName:         cc.StructName,
		ShortStructName:    strings.ToLower(cc.StructName[0:1]),
		StructFields:       fields,
		TableIndexes:       uniqueStrings(tableIndexes),
		TableUniques:       uniqueStrings(tableUniques),
		EnableFieldComment: cc.EnableFieldComment,
		EnableSqlNull:      cc.EnableSqlNull,
		EnableGureguNull:   cc.EnableGureguNull,
		EnableTableName:    cc.EnableGormTag || cc.EnableXormTag || cc.EnableBeegoTag || cc.EnableGoroseTag,
		EnableTableIndex:   cc.EnableBeegoTag && len(tableIndexes) != 0,
		EnableTableUnique:  cc.EnableBeegoTag && len(tableUniques) != 0,
	})
	if err != nil {
		fmt.Println("execute template err:", err)
		return "", err
	}

	code, err := format.Source(buffer.Bytes())
	if err != nil {
		fmt.Println("go fmt err:", err, buffer.Bytes())
		return "", err
	}

	return string(code[:len(code)-1]), nil
}

// generateTag generates the tag string by column information and tag name.
func generateTag(ci *ColumnInfo, tag string) string {
	buffer := &bytes.Buffer{}
	err := generator.ExecuteTemplate(buffer, tag, ci)
	if err != nil {
		fmt.Println("execute template err:", err)
		return ""
	}

	return strings.TrimSpace(buffer.String())
}

// uniqueStrings returns the unique string slice.
func uniqueStrings(slice []string) []string {
	var result []string
	uniqueMap := make(map[string]bool)

	for _, value := range slice {
		if _, ok := uniqueMap[value]; !ok {
			uniqueMap[value] = true
			result = append(result, value)
		}
	}

	return result
}
