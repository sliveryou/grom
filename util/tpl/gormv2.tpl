gorm:"
{{- if .IsPrimaryKey }}primaryKey;{{ end -}}
{{ if .IsAutoIncrement }}autoIncrement;{{ end }}column:{{ .Name }}{{ if not .IsPrimaryKey }};type:{{ .Type }}{{ end }}
{{- if or .IsNullable .IsPrimaryKey | not }};not null{{ end -}}
{{- range $i, $v := .Indexes }}
    {{- if eq $i 0 }};index:{{ $v.Name }}{{ else }},{{ $v.Name }}{{ end }}{{ end -}}
{{- range $i, $v := .UniqueIndexes }}
    {{- if eq $i 0 }};uniqueIndex:{{ $v.Name }}{{ else }},{{ $v.Name }}{{ end }}{{ end -}}
{{- if .Default }};default:{{ .Default }}{{ end -}}
{{- if .Comment }};comment:{{ .Comment }}{{ end -}}
"