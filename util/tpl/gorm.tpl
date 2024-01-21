gorm:"
{{- if .IsPrimaryKey }}primary_key;{{ end -}}
column:{{ .Name }};type:{{ .Type }}{{ if .IsAutoIncrement }} auto_increment{{ end }}
{{- if or .IsNullable .IsPrimaryKey | not }};not null{{ end -}}
{{- range $i, $v := .Indexes }}
    {{- if eq $i 0 }};index:{{ $v.Name }}{{ else }},{{ $v.Name }}{{ end }}{{ end -}}
{{- range $i, $v := .UniqueIndexes }}
    {{- if eq $i 0 }};unique_index:{{ $v.Name }}{{ else }},{{ $v.Name }}{{ end }}{{ end -}}
{{- if .Default }};default:'{{ .Default }}'{{ end -}}
{{- if .Comment }};comment:'{{ .Comment }}'{{ end -}}
"