orm:"
{{- if .IsPrimaryKey }}pk;{{ end -}}
{{- if .IsAutoIncrement }}auto;{{ end -}}
column({{ .Name }}){{ getBeegoType . }}
{{- if .IsNullable }};null{{ end -}}
{{- if .Default }};default({{ .Default }}){{ end -}}
{{- if .Comment }};description({{ .Comment }}){{ end -}}
"