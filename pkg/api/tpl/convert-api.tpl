// To{{ .StructName }} {{ .TableComment }}转化
func To{{ .StructName }}(src *pb.{{ .StructName }}) {{ if .IfInfo }}({{ .StructName }}, error){{ else }}{{ .StructName }}{{ end }} {
	var dst {{ .StructName }}
	if src != nil {
		dst = {{ .StructName }}{
			{{ .ConvertInfo }}
		}
		{{ .IfInfo -}}
	}

	return dst{{ if .IfInfo }}, nil{{ end }}
}
