{{ if .IsPointer -}}
{{ if .IsStringType -}}
if in.{{ .Name }} != nil {
	{{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Like("%" + *in.{{ .Name }} + "%"))
}
{{ else if .IsNumberType -}}
if in.{{ .Name }} != nil {
	{{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Eq(*in.{{ .Name }}))
}
{{ else if .IsTimeType -}}
if in.{{ .Name }} != nil {
	{{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Eq(time.UnixMilli(*in.{{ .Name }})))
}
{{ else if .IsBoolType -}}
if in.{{ .Name }} != nil {
    {{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Is(*in.{{ .Name }}))
}
{{ end -}}
{{ else -}}
{{ if .IsStringType -}}
if in.{{ .Name }} != "" {
	{{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Like("%" + in.{{ .Name }} + "%"))
}
{{ else if .IsNumberType -}}
if in.{{ .Name }} != 0 {
	{{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Eq(in.{{ .Name }}))
}
{{ else if .IsTimeType -}}
if in.{{ .Name }} != 0 {
	{{ .SmallStructName }}q = {{ .SmallStructName }}q.Where({{ .SmallStructName }}.{{ .ReplaceName }}.Eq(time.UnixMilli(in.{{ .Name }})))
}
{{ end -}}
{{ end -}}
