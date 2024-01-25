{{ if .IsNullable -}}
if in.{{ .MemberName }} != nil && *in.{{ .MemberName }} != {{ .ObjectName }}.{{ .ObjectMemberName }} {
    updateMap["{{ .MemberRawName }}"] = *in.{{ .MemberName }}
    {{ .ObjectName }}.{{ .ObjectMemberName }} = *in.{{ .MemberName }}
}
{{ else -}}
if in.{{ .MemberName }} != {{ .ObjectName }}.{{ .ObjectMemberName }} {
    updateMap["{{ .MemberRawName }}"] = in.{{ .MemberName }}
    {{ .ObjectName }}.{{ .ObjectMemberName }} = in.{{ .MemberName }}
}
{{ end -}}
