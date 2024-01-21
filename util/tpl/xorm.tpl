xorm:"
{{- if .IsPrimaryKey }}pk {{ end -}}
{{- if .IsAutoIncrement }}autoincr {{ end -}}
{{ .Type }} '{{ .Name }}'
{{- if or .IsNullable .IsPrimaryKey | not }} notnull{{ end -}}
{{- range $i, $v := .Indexes }} index({{ $v.Name }}){{ end -}}
{{- range $i, $v := .UniqueIndexes }} unique({{ $v.Name }}){{ end -}}
{{- if .Default }} default('{{ .Default }}'){{ end -}}
{{- if .Comment }} comment('{{ .Comment }}'){{ end -}}
"