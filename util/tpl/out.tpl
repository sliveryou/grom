package {{.PackageName}}

{{ if .NeedImport }}
import (
	{{ if .EnableGureguNull }}
		"gopkg.in/guregu/null.v4"
	{{ else if .EnableSqlNull }}
		"database/sql"
	{{ else if .EnableGoTime }}
		"time"
	{{ end }}
)
{{ end }}

// {{ .StructName }} {{ .TableComment }}
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