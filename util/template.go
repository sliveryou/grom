package util

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var generator *template.Template

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

func init() {
	var err error
	generator, err = template.New("gorm").Parse(gormTemplate)
	if err != nil {
		fmt.Println("parse gorm template err:", err)
	}
	generator, err = generator.New("xorm").Parse(xormTemplate)
	if err != nil {
		fmt.Println("parse xorm template err:", err)
	}
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
