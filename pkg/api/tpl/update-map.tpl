{{ if .HasDefault -}}
if in.{{ .MemberName }} != nil {{ if not .IsDataTypeJSON }}&& ({{ .ObjectName }}.{{ .ObjectMemberName }} == nil || *in.{{ .MemberName }} != *{{ .ObjectName }}.{{ .ObjectMemberName }}){{ end }} {
    updateMap["{{ .MemberRawName }}"] = {{ if not .IsDataTypeJSON }}*{{ end }}in.{{ .MemberName }}
    {{ .ObjectName }}.{{ .ObjectMemberName }} = in.{{ .MemberName }}
}
{{ else if .IsNullable -}}
{{ if .IsTimeField -}}
if in.{{ .MemberName }} != nil {
    {{ .MemberLowerCamelName }} := time.UnixMilli(*in.{{ .MemberName }})
    if {{ .ObjectName }}.{{ .ObjectMemberName }} == nil || !{{ .ObjectName }}.{{ .ObjectMemberName }}.Equal({{ .MemberLowerCamelName }}) {
        updateMap["{{ .MemberRawName }}"] = &{{ .MemberLowerCamelName }}
        {{ .ObjectName }}.{{ .ObjectMemberName }} = &{{ .MemberLowerCamelName }}
    }
}
{{ else -}}
if in.{{ .MemberName }} != nil {{ if not .IsDataTypeJSON }}&& *in.{{ .MemberName }} != {{ .ObjectName }}.{{ .ObjectMemberName }}{{ end }} {
    updateMap["{{ .MemberRawName }}"] = {{ if not .IsDataTypeJSON }}*{{ end }}in.{{ .MemberName }}
    {{ .ObjectName }}.{{ .ObjectMemberName }} = {{ if not .IsDataTypeJSON }}*{{ end }}in.{{ .MemberName }}
}
{{ end -}}
{{ else -}}
{{ if .IsTimeField -}}
if new{{ .MemberName }} := time.UnixMilli(in.{{ .MemberName }}); {{ .ObjectName }}.{{ .ObjectMemberName }} == nil || !{{ .ObjectName }}.{{ .ObjectMemberName }}.Equal(new{{ .MemberName }}) {
    updateMap["{{ .MemberRawName }}"] = &new{{ .MemberName }}
    {{ .ObjectName }}.{{ .ObjectMemberName }} = &new{{ .MemberName }}
}
{{ else -}}
if {{ if .IsPointer }}*{{ end }}in.{{ .MemberName }} != {{ .ObjectName }}.{{ .ObjectMemberName }} {
    updateMap["{{ .MemberRawName }}"] = {{ if .IsPointer }}*{{ end }}in.{{ .MemberName }}
    {{ .ObjectName }}.{{ .ObjectMemberName }} = {{ if .IsPointer }}*{{ end }}in.{{ .MemberName }}
}
{{ end -}}
{{ end -}}
